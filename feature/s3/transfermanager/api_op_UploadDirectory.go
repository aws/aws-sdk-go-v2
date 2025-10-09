package transfermanager

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// UploadDirectoryInput represents a request to the UploadDirectory() call
type UploadDirectoryInput struct {
	// Bucket where objects are uploaded to
	Bucket string

	// The source directory to upload
	Source string

	// Whether to follow symbolic links when traversing the file tree.
	FollowSymbolicLinks bool

	// Whether to recursively upload directories. If set to false by
	// default, only top level files under source folder will be uplaoded;
	// otherwise all files under subfolders will be uploaded
	Recursive bool

	// The S3 key prefix to use for each object. If not provided, files
	// will be uploaded to the root of the bucket
	KeyPrefix string

	// A callback func to allow users to filter out unwanted files
	// according to bool returned from the function
	Filter FileFilter

	// A callback function to allow customers to update individual
	// PutObjectInput that the S3 Transfer Manager generates.
	Callback PutRequestCallback

	// The s3 delimeter contatenating each object key based on local file separator
	// and file's relative path. If a non-defualt delimiter is used, it can not be
	// included in any subfolders or files, which will cause error otherwise
	S3Delimiter string
}

// FileFilter is the callback to allow users to filter out unwanted files.
// It is invoked for each file.
type FileFilter interface {
	// FilterFile take the file path and decides if the file
	// should be uploaded
	FilterFile(filePath string) bool
}

// PutRequestCallback is the callback mechanism to allow customers to update
// individual PutObjectInput that the S3 Transfer Manager generates
type PutRequestCallback interface {
	// UpdateRequest preprocesses each PutObjectInput as customized
	UpdateRequest(*UploadObjectInput)
}

// UploadDirectoryOutput represents a response from the UploadDirectory() call
type UploadDirectoryOutput struct {
	// Total number of objects successfully uploaded
	ObjectsUploaded int
}

// UploadDirectory traverses a local directory recursively/non-recursively and intelligently
// uploads all valid files to S3 in parallel across multiple goroutines. You can configure
// the concurrency, valid file filtering and object key naming through the Options and input parameters.
//
// Additional functional options can be provided to configure the individual directory
// upload. These options are copies of the original Options instance, the client of which UploadDirectory is called from.
// Modifying the options will not impact the original Client and Options instance.
func (c *Client) UploadDirectory(ctx context.Context, input *UploadDirectoryInput, opts ...func(*Options)) (*UploadDirectoryOutput, error) {
	fileInfo, err := os.Stat(input.Source)
	if err != nil {
		return nil, fmt.Errorf("error when getting source info: %v", err)
	}
	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("the source path %s doesn't point to a valid directory", input.Source)
	}

	i := directoryUploader{c: c, in: input, options: c.options.Copy()}
	for _, opt := range opts {
		opt(&i.options)
	}

	return i.uploadDirectory(ctx)
}

type directoryUploader struct {
	c       *Client
	options Options
	in      *UploadDirectoryInput

	filesUploaded int
	traversed     map[string]interface{}

	err error

	mu           sync.Mutex
	wg           sync.WaitGroup
	progressOnce sync.Once

	emitter *directoryObjectsProgressEmitter
}

func (u *directoryUploader) uploadDirectory(ctx context.Context) (*UploadDirectoryOutput, error) {
	u.init()
	ch := make(chan fileEntry)

	for i := 0; i < u.options.DirectoryConcurrency; i++ {
		u.wg.Add(1)
		go u.uploadFile(ctx, ch)
	}

	if u.in.Recursive {
		u.traverse(u.in.Source, u.in.KeyPrefix, ch)
	} else {
		files, err := u.traverseFolder(u.in.Source)
		if err != nil {
			return nil, err
		}

		for _, f := range files {
			if u.getErr() != nil {
				break
			}

			path := filepath.Join(u.in.Source, f)
			absPath, err := u.getAbsPath(path)
			if err != nil {
				u.setErr(fmt.Errorf("error when getting abs path of file %s: %v", path, err))
				break
			} else if absPath == "" {
				continue
			}

			fileInfo, err := os.Lstat(absPath)
			if err != nil {
				u.setErr(fmt.Errorf("error when stating abs path %s: %v", absPath, err))
				break
			}
			if fileInfo.IsDir() {
				continue
			}

			if u.in.Filter != nil && !u.in.Filter.FilterFile(path) {
				continue
			}
			if u.in.S3Delimiter != "/" && strings.Contains(f, u.in.S3Delimiter) {
				u.setErr(fmt.Errorf("file %s contains delimiter %s", f, u.in.S3Delimiter))
				break
			}
			if u.in.KeyPrefix == "" {
				ch <- fileEntry{f, absPath}
			} else {
				ch <- fileEntry{u.in.KeyPrefix + u.in.S3Delimiter + f, absPath}
			}
		}
	}
	close(ch)
	u.wg.Wait()

	if u.err != nil {
		u.emitter.Failed(ctx, u.in, u.err)
		return nil, u.err
	}

	out := &UploadDirectoryOutput{
		ObjectsUploaded: u.filesUploaded,
	}
	u.emitter.Complete(ctx, out)
	return out, nil
}

func (u *directoryUploader) init() {
	if u.in.S3Delimiter == "" {
		u.in.S3Delimiter = "/"
	}

	u.traversed = make(map[string]interface{})

	u.emitter = &directoryObjectsProgressEmitter{
		Listeners: u.options.DirectoryProgressListeners,
	}
}

type fileEntry struct {
	key  string
	path string
}

// traverse recursively visits each folder and sends each
// valid file's request to worker goroutines
func (u *directoryUploader) traverse(path, keyPrefix string, ch chan fileEntry) {
	if u.getErr() != nil {
		return
	}

	absPath, err := u.getAbsPath(path)
	if err != nil {
		u.setErr(fmt.Errorf("error when getting abs path of file %s: %v", path, err))
		return
	} else if absPath == "" {
		return
	}

	var key string
	if path == u.in.Source {
		key = keyPrefix
	} else if keyPrefix == "" {
		key = filepath.Base(path)
	} else {
		key = keyPrefix + u.in.S3Delimiter + filepath.Base(path)
	}
	fileInfo, err := os.Lstat(absPath)
	if err != nil {
		u.setErr(fmt.Errorf("error when stating file %s: %v", absPath, err))
		return
	}
	if fileInfo.IsDir() {
		subFiles, err := u.traverseFolder(absPath)
		if err != nil {
			u.setErr(fmt.Errorf("error when traversing folder %s: %v", absPath, err))
			return
		}
		for _, f := range subFiles {
			if d := u.in.S3Delimiter; d != "/" && strings.Contains(f, d) {
				u.setErr(fmt.Errorf("file %s contains delimiter %s", f, d))
				return
			}
			u.traverse(filepath.Join(path, f), key, ch)
		}
	} else {
		if u.in.Filter != nil && !u.in.Filter.FilterFile(path) {
			return
		}
		if u.in.S3Delimiter != "/" {
			if n, d := filepath.Base(path), u.in.S3Delimiter; strings.Contains(n, d) {
				u.setErr(fmt.Errorf("file %s contains delimiter %s", n, d))
				return
			}
		}
		ch <- fileEntry{key, absPath}
	}
}

// getAbsPath resolves a path's desination absolute path with deduplication
// in case any symlink causes traverse loop or repeated upload
func (u *directoryUploader) getAbsPath(path string) (string, error) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return "", fmt.Errorf("error when stating info of file %s: %v", path, err)
	}

	if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		if !u.in.FollowSymbolicLinks {
			return "", nil
		}
		path, err = u.traverseSymlink(path)
		if err != nil {
			return "", err
		}
	}
	if u.traversed[path] != nil {
		return "", fmt.Errorf("traversed duplicate path %s", path)
	}
	u.traversed[path] = struct{}{}

	return path, nil
}

// traverseFolder returns subfiles at this level
func (u *directoryUploader) traverseFolder(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return []string{}, err
	}
	subFiles, err := f.ReadDir(0)
	if err != nil {
		return []string{}, err
	}

	files := []string{}
	for _, v := range subFiles {
		files = append(files, v.Name())
	}

	return files, nil
}

func (u *directoryUploader) traverseSymlink(path string) (string, error) {
	originPath := path
	for {
		dst, err := os.Readlink(path)
		if err != nil {
			return "", fmt.Errorf("error when reading symlink of %s: %v", originPath, err)
		}
		if filepath.IsAbs(dst) {
			path = dst
		} else {
			path = filepath.Join(filepath.Dir(path), dst)
		}
		if u.traversed[path] != nil {
			return "", fmt.Errorf("traversed duplicate path: %s", path)
		}
		fileInfo, err := os.Lstat(path)
		if err != nil {
			return "", fmt.Errorf("error when stating linked path %s: %v", path, err)
		}
		if fileInfo.Mode()&os.ModeSymlink != os.ModeSymlink {
			return path, nil
		}
		u.traversed[path] = struct{}{}
	}
}

func (u *directoryUploader) uploadFile(ctx context.Context, ch chan fileEntry) {
	defer u.wg.Done()

	for {
		data, ok := <-ch
		if !ok {
			break
		}

		select {
		case <-ctx.Done():
			u.setErr(fmt.Errorf("context error: %v", ctx.Err()))
			continue
		default:
		}

		if u.getErr() != nil {
			continue
		}
		f, err := os.Open(data.path)
		if err != nil {
			u.setErr(fmt.Errorf("error when opening file %s: %v", data.path, err))
			continue
		}
		input := &UploadObjectInput{
			Bucket: u.in.Bucket,
			Key:    data.key,
			Body:   f,
		}
		if u.in.Callback != nil {
			u.in.Callback.UpdateRequest(input)
		}
		out, err := u.c.UploadObject(ctx, input)
		if err != nil {
			u.setErr(fmt.Errorf("error when uploading file %s: %v", data.path, err))
			continue
		}

		u.progressOnce.Do(func() {
			u.emitter.Start(ctx, u.in)
		})
		u.incrFilesUploaded(1)
		u.emitter.ObjectsTransferred(ctx, out.ContentLength)
	}
}

func (u *directoryUploader) incrFilesUploaded(n int) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.filesUploaded += n
}

func (u *directoryUploader) setErr(err error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.err = err
}

func (u *directoryUploader) getErr() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	return u.err
}
