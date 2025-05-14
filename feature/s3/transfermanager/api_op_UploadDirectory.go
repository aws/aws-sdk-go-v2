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
		files, _, err := u.traverseFolder(u.in.Source)
		if err != nil {
			return nil, err
		}

		for _, f := range files {
			filePath := filepath.Join(u.in.Source, f)
			if u.in.Filter != nil && !u.in.Filter.FilterFile(filePath) {
				continue
			}
			key, err := u.mapKeyFromPath(f, u.in.KeyPrefix)
			if err != nil {
				u.setErr(err)
				break
			}
			ch <- fileChunk{key: key, path: filePath}
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
	if u.in.KeyPrefix != "" && !strings.HasSuffix(u.in.KeyPrefix, u.in.S3Delimiter) {
		u.in.KeyPrefix = u.in.KeyPrefix + u.in.S3Delimiter
	}
}

type fileChunk struct {
	key  string
	path string
}

source b/
a/b/c -> d/e/ folder
	 d/e/f	file key a/c/f
	 d/e/g/ folder 
	 d/e/g/h file key a/b/c/g/h
	
	

func (u *directoryUploader) traverse(folderPath, keyPrefix string, ch chan fileChunk) {
	if u.getErr() != nil {
		return
	}
	files, directories, err := u.traverseFolder(folderPath)
	if err != nil {
		u.setErr(err)
		return
	}

	for _, f := range files {
		filePath := filepath.Join(folderPath, f)
		if u.in.Filter != nil && !u.in.Filter.FilterFile(filePath) {
			continue
		}
		key, err := u.mapKeyFromPath(f, keyPrefix)
		if err != nil {
			u.setErr(err)
			break
		}
		ch <- fileChunk{key: key, path: filePath}
	}

	for _, d := range directories {
		u.traverse(filepath.Join(folderPath, d), keyPrefix+d+u.in.S3Delimiter, ch)
	}
}

func (u *directoryUploader) traverseFolder(path string) (files, directories []string, err error) {
	f, e := os.Open(path)
	if e != nil {
		err = e
		return
	}
	subFiles, e := f.Readdir(0)
	if e != nil {
		err = e
		return
	}

	for _, v := range subFiles {
		if v.IsDir() {
			directories = append(directories, v.Name())
		} else {
			files = append(files, v.Name())
		}
	}

	return
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
