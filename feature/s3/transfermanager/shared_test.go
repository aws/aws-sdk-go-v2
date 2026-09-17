package transfermanager

import "fmt"

var buf20MB = make([]byte, 1024*1024*20)
var buf2MB = make([]byte, 1024*1024*2)
var buf40MB = make([]byte, 1024*1024*40)
var buf80MB = make([]byte, 1024*1024*80)

// describeBytesDiff summarizes how two byte slices differ, for assertions on
// payloads too large to print. Misplaced writes preserve length, so the offset
// of the first difference is reported as well as the lengths.
func describeBytesDiff(expect, actual []byte) string {
	if len(expect) != len(actual) {
		return fmt.Sprintf("expect %d bytes, got %d bytes", len(expect), len(actual))
	}
	for i := range expect {
		if expect[i] != actual[i] {
			return fmt.Sprintf("both %d bytes, first difference at offset %d: expect %q, got %q",
				len(expect), i, expect[i], actual[i])
		}
	}
	return "no difference"
}
