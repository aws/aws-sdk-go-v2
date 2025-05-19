package transfermanager

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

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
	// and file's relative path
	S3Delimiter string
}

type FileFilter interface {
	FilterFile(filePath string) bool
}

type PutRequestCallback interface {
	UpdateRequest(*PutObjectInput)
}

type UploadDirectoryOutput struct {
	// Total number of objects successfully uploaded
	ObjectsUploaded int

	// Total number of objects failed during upload
	ObjectsFailed int
}

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
	filesFailed   int
	traversed     map[string]interface{}

	err error

	mu sync.Mutex
	wg sync.WaitGroup
}

func (u *directoryUploader) uploadDirectory(ctx context.Context) (*UploadDirectoryOutput, error) {
	u.init()
	ch := make(chan fileChunk)

	for i := 0; i < u.options.Concurrency; i++ {
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
			path := filepath.Join(u.in.Source, f)
			if u.in.Filter != nil && !u.in.Filter.FilterFile(path) {
				continue
			}
			if u.in.S3Delimiter != "/" && strings.Contains(f, u.in.S3Delimiter) {
				return nil, fmt.Errorf("file %s contains delimiter %s", f, u.in.S3Delimiter)
			}
			fileInfo, err := os.Lstat(path)
			if err != nil {
				return nil, err
			}
			if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
				if !u.in.FollowSymbolicLinks {
					continue
				}
				path, err = u.traverseSymlink(path)
				if err != nil {
					return nil, err
				}
				fileInfo, err = os.Lstat(path)
				if err != nil {
					return nil, err
				}
			}

			if fileInfo.IsDir() {
				continue
			}
			if u.traversed[path] != nil {
				return nil, fmt.Errorf("traversed duplicate path %s", path)
			}
			u.traversed[path] = struct{}{}
			ch <- fileChunk{u.in.KeyPrefix + u.in.S3Delimiter + f, path}
		}
	}
	close(ch)
	u.wg.Wait()

	if u.err != nil {
		return nil, u.err
	}
	return &UploadDirectoryOutput{
		ObjectsUploaded: u.filesUploaded,
		ObjectsFailed:   u.filesFailed,
	}, nil
}

func (u *directoryUploader) init() {
	if u.in.S3Delimiter == "" {
		u.in.S3Delimiter = "/"
	}
}

type fileChunk struct {
	key  string
	path string
}

func (u *directoryUploader) traverse(path, keyPrefix string, ch chan fileChunk) {
	if u.getErr() != nil {
		return
	}

	fileInfo, err := os.Lstat(path)
	if err != nil {
		u.setErr(err)
		return
	}
	absPath := path
	if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		if !u.in.FollowSymbolicLinks {
			return
		}
		absPath, err = u.traverseSymlink(absPath)
		if err != nil {
			u.setErr(err)
			return
		}
	}

	if u.traversed[absPath] != nil {
		u.setErr(fmt.Errorf("traversed duplicate path: %s", absPath))
		return
	}
	u.traversed[absPath] = struct{}{}

	fileInfo, err = os.Lstat(absPath)
	if err != nil {
		u.setErr(err)
		return
	}

	if path == u.in.Source {
		key := keyPrefix
	} else {
		key := keyPrefix + u.in.S3Delimiter + filepath.Base(path)
	}
	if fileInfo.IsDir() {
		subFiles, err := u.traverseFolder(absPath)
		if err != nil {
			u.setErr(err)
			return
		}
		for _, f := range subFiles {
			u.traverse(path+f, key, ch)
		}
	} else {
		if u.in.Filter != nil && !u.in.Filter.FilterFile(absPath) {
			return
		}
		if u.in.S3Delimiter != "/" {
			if n, d := fileInfo.Name(), u.in.S3Delimiter; strings.Contains(n, d) {
				u.setErr(fmt.Errorf("file %s contains delimiter %s", n, d))
				return
			}
		}
		ch <- fileChunk{key, absPath}
	}
}

func (u *directoryUploader) traverseFolder(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return []string{}, err
	}
	subFiles, err := f.Readdir(0)
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
	for {
		dst, err := os.Readlink(path)
		if err != nil {
			return "", err
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
			return "", err
		}
		if fileInfo.Mode()&os.ModeSymlink != os.ModeSymlink {
			return path, nil
		}
		u.traversed[path] = struct{}{}
	}
}

func (u *directoryUploader) uploadFile(ctx context.Context, ch chan fileChunk) {
	defer u.wg.Done()

	for {
		data, ok := <-ch
		if !ok {
			break
		}
		if u.getErr() != nil {
			continue
		}
		f, err := os.Open(data.path)
		if err != nil {
			u.setErr(err)
		} else {
			input := &PutObjectInput{
				Bucket: u.in.Bucket,
				Key:    data.key,
				Body:   f,
			}
			if u.in.Callback != nil {
				u.in.Callback.UpdateRequest(input)
			}
			_, err := u.c.PutObject(ctx, input)
			if err != nil {
				u.setErr(err)
				u.incrFilesFailed(1)
			} else {
				u.incrFilesUploaded(1)
			}
		}
	}
}

func (u *directoryUploader) mapKeyFromPath(fileName, keyPrefix string) (string, error) {
	if u.in.S3Delimiter != "/" && strings.Contains(fileName, u.in.S3Delimiter) {
		return "", fmt.Errorf("file %s contains the delimiter %s", fileName, u.in.S3Delimiter)
	}
	return keyPrefix + fileName, nil
}

func (u *directoryUploader) incrFilesUploaded(n int) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.filesUploaded += n
}

func (u *directoryUploader) incrFilesFailed(n int) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.filesFailed += n
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
