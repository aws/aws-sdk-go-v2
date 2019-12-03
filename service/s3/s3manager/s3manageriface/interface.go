// Package s3manageriface provides an interface for the s3manager package
package s3manageriface

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// DownloaderAPI is the interface type for s3manager.Downloader.
type DownloaderAPI interface {
	Download(io.WriterAt, *types.GetObjectInput, ...func(*s3manager.Downloader)) (int64, error)
	DownloadWithContext(context.Context, io.WriterAt, *types.GetObjectInput, ...func(*s3manager.Downloader)) (int64, error)
	DownloadWithIterator(context.Context, s3manager.BatchDownloadIterator, ...func(*s3manager.Downloader)) error
}

var _ DownloaderAPI = (*s3manager.Downloader)(nil)

// UploaderAPI is the interface type for s3manager.Uploader.
type UploaderAPI interface {
	Upload(*s3manager.UploadInput, ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
	UploadWithContext(context.Context, *s3manager.UploadInput, ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
	UploadWithIterator(context.Context, s3manager.BatchUploadIterator, ...func(*s3manager.Uploader)) error
}

var _ UploaderAPI = (*s3manager.Uploader)(nil)
