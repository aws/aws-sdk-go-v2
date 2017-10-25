// +build !appengine,!plan9

package aws_test

import (
	"net"
	"os"
	"syscall"
)

var stubConnectionResetError = &net.OpError{Err: &os.SyscallError{Syscall: "read", Err: syscall.ECONNRESET}}
