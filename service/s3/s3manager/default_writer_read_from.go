// +build !windows,disabled

package s3manager

func defaultDownloadBufferProvider() WriterReadFromProvider {
	return nil
}
