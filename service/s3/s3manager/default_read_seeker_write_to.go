// +build !windows,disabled

package s3manager

func defaultUploadBufferProvider() ReadSeekerWriteToProvider {
	return nil
}
