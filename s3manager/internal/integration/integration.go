package integration

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"

	smithyrand "github.com/awslabs/smithy-go/rand"
)

var uuid = smithyrand.NewUUID(rand.Reader)

// MustUUID returns an UUID string or panics
func MustUUID() string {
	uuid, err := uuid.GetUUID()
	if err != nil {
		panic(err)
	}
	return uuid
}

// CreateFileOfSize will return an *os.File that is of size bytes
func CreateFileOfSize(dir string, size int64) (*os.File, error) {
	file, err := ioutil.TempFile(dir, "s3integration")
	if err != nil {
		return nil, err
	}

	err = file.Truncate(size)
	if err != nil {
		file.Close()
		os.Remove(file.Name())
		return nil, err
	}

	return file, nil
}

// SizeToName returns a human-readable string for the given size bytes
func SizeToName(size int) string {
	units := []string{"B", "KB", "MB", "GB"}
	i := 0
	for size >= 1024 {
		size /= 1024
		i++
	}

	if i > len(units)-1 {
		i = len(units) - 1
	}

	return fmt.Sprintf("%d%s", size, units[i])
}
