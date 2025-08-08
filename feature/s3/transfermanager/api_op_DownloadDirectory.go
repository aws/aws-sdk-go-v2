package transfermanager

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// DownloadDirectoryInput represents a request to the DownloadDirectory() call
type DownloadDirectoryInput struct {
	// Bucket where objects are downloaded from
	Bucket string

	// The destination directory to download
	Destination string

	// The S3 key prefix to use for listing objects. If not provided,
	// all objects under a bucket will be retrieved
	KeyPrefix string

	// The s3 delimiter used to convert keyname to local filepath if it
	// is different from local file separator
	S3Delimiter string

	// A callback func to allow users to fileter out unwanted objects
	// according to bool returned from the function
	Filter ObjectFilter

	// A callback function to allow customers to update individual
	// GetObjectInput that the S3 Transfer Manager generates
	Callback GetRequestCallback
}

// ObjectFilter is the callback to allow users to filter out unwanted objects.
// It is invoked for each object listed.
type ObjectFilter interface {
	// FilterObject take the Object struct and decides if the
	// object should be downloaded
	FilterObject(s3types.Object) bool
}

// GetRequestCallback is the callback mechanism to allow customers to update
// individual GetObjectInput that the S3 Transfer Manager generates
type GetRequestCallback interface {
	// UpdateRequest preprocesses each GetObjectInput as customized
	UpdateRequest(*GetObjectInput)
}

// DownloadDirectoryOutput represents a response from the DownloadDirectory() call
type DownloadDirectoryOutput struct {
	// Total number of objects successfully downloaded
	ObjectsDownloaded int
}

type objectEntry struct {
	key  string
	path string
}

// DownloadDirectory traverses a s3 bucket and intelligently downloads all valid objects
// to local directory in parallel across multiple goroutines. You can configure the concurrency,
// valid object filtering and hierarchical file naming through the Options and input parameters.
//
// Additional functional options can be provided to configure the individual directory
// download. These options are copies of the original Options instance, the client of which DownloadDirectory is called from.
// Modifying the options will not impact the original Client and Options instance.
func (c *Client) DownloadDirectory(ctx context.Context, input *DownloadDirectoryInput, opts ...func(*Options)) (*DownloadDirectoryOutput, error) {
	fileInfo, err := os.Stat(input.Destination)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("error when getting destination folder info: %v", err)
		}
	} else {
		if !fileInfo.IsDir() {
			return nil, fmt.Errorf("the destination path %s doesn't point to a valid directory", input.Destination)
		}
	}

	i := directoryDownloader{c: c, in: input, options: c.options.Copy()}
	for _, opt := range opts {
		opt(&i.options)
	}

	return i.downloadDirectory(ctx)
}

type directoryDownloader struct {
	c       *Client
	options Options
	in      *DownloadDirectoryInput

	objectsDownloaded int

	err error

	mu sync.Mutex
	wg sync.WaitGroup
}

func (d *directoryDownloader) downloadDirectory(ctx context.Context) (*DownloadDirectoryOutput, error) {
	d.init()
	ch := make(chan objectEntry)

	for i := 0; i < d.options.DirectoryConcurrency; i++ {
		d.wg.Add(1)
		go d.downloadObject(ctx, ch)
	}

	isTruncated := true
	continuationToken := ""
	for isTruncated {
		if d.getErr() != nil {
			break
		}
		listOutput, err := d.options.S3.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket:            aws.String(d.in.Bucket),
			Prefix:            nzstring(d.in.KeyPrefix),
			ContinuationToken: nzstring(continuationToken),
		})
		if err != nil {
			d.setErr(fmt.Errorf("error when listing objects %v", err))
			break
		}

		for _, o := range listOutput.Contents {
			key := aws.ToString(o.Key)
			if strings.HasSuffix(key, "/") || strings.HasSuffix(key, d.in.S3Delimiter) {
				continue // skip folder object
			}
			if d.in.Filter != nil && !d.in.Filter.FilterObject(o) {
				continue
			}
			path, err := d.getLocalPath(key)
			if err != nil {
				d.setErr(fmt.Errorf("error when resolving local path for object %s, %v", key, err))
				break
			}
			ch <- objectEntry{key, path}
		}

		continuationToken = aws.ToString(listOutput.NextContinuationToken)
		isTruncated = aws.ToBool(listOutput.IsTruncated)
	}

	close(ch)
	d.wg.Wait()

	if d.err != nil {
		return nil, d.err
	}

	return &DownloadDirectoryOutput{
		ObjectsDownloaded: d.objectsDownloaded,
	}, nil
}

func (d *directoryDownloader) init() {
	if d.in.S3Delimiter == "" {
		d.in.S3Delimiter = "/"
	}
}

func (d *directoryDownloader) getLocalPath(key string) (string, error) {
	keyprefix := d.in.KeyPrefix
	if keyprefix != "" && !strings.HasSuffix(keyprefix, d.in.S3Delimiter) {
		keyprefix = keyprefix + d.in.S3Delimiter
	}
	path := filepath.Join(d.in.Destination, strings.ReplaceAll(strings.TrimPrefix(key, keyprefix), d.in.S3Delimiter, string(os.PathSeparator)))
	relPath, err := filepath.Rel(d.in.Destination, path)
	if err != nil {
		return "", err
	}
	if relPath == "." || strings.Contains(relPath, "../") {
		return "", fmt.Errorf("resolved local path %s is outside of destination %s", path, d.in.Destination)
	}

	return path, nil
}

func (d *directoryDownloader) downloadObject(ctx context.Context, ch chan objectEntry) {
	defer d.wg.Done()
	for {
		data, ok := <-ch
		if !ok {
			break
		}
		if d.getErr() != nil {
			break
		}

		input := &GetObjectInput{
			Bucket: d.in.Bucket,
			Key:    data.key,
		}
		if d.in.Callback != nil {
			d.in.Callback.UpdateRequest(input)
		}
		out, err := d.c.GetObject(ctx, input)
		if err != nil {
			d.setErr(fmt.Errorf("error when downloading object %s: %v", data.key, err))
			break
		}

		err = os.MkdirAll(filepath.Dir(data.path), os.ModePerm)
		if err != nil {
			d.setErr(fmt.Errorf("error when creating directory for file %s: %v", data.path, err))
			break
		}
		file, err := os.Create(data.path)
		if err != nil {
			d.setErr(fmt.Errorf("error when creating file %s: %v", data.path, err))
			break
		}
		_, err = io.Copy(file, out.Body)
		if err != nil {
			d.setErr(fmt.Errorf("error when writing to local file %s: %v", data.path, err))
			os.Remove(data.path)
			break
		}
		d.incrObjectsDownloaded(1)
	}
}

func (d *directoryDownloader) incrObjectsDownloaded(n int) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.objectsDownloaded += n
}

func (d *directoryDownloader) setErr(err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.err = err
}

func (d *directoryDownloader) getErr() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.err
}
