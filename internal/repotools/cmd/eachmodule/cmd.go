package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

// Work provides a pending job to be done.
type Work struct {
	Path string
	Cmd  string
}

// WorkLog provides the result of a job.
type WorkLog struct {
	Path, Cmd string
	Err       error
	Output    io.Reader
}

// CommandWorker provides a consumer of work jobs and posts results to the
// worklog.
func CommandWorker(ctx context.Context, jobs <-chan Work, results chan<- WorkLog, streamOut io.Writer) {
	for {
		var result WorkLog

		select {
		case <-ctx.Done():
			return
		case w, ok := <-jobs:
			if !ok {
				return
			}

			outBuffer := bytes.NewBuffer(nil)
			outWriter := io.Writer(outBuffer)

			if streamOut != nil {
				outWriter = io.MultiWriter(outWriter, streamOut)
			}

			result.Path = w.Path
			result.Cmd = w.Cmd

			cmd, err := NewCommand(ctx, outWriter, outWriter, w.Path, w.Cmd)
			if err != nil {
				result.Err = fmt.Errorf("failed to build command, %w", err)
				break
			}

			if err := cmd.Run(); err != nil {
				result.Err = fmt.Errorf("failed to run command, %v", err)
			}

			if streamOut == nil {
				outReader := bytes.NewReader(outBuffer.Bytes())
				result.Output = outReader
			}
		}

		select {
		case <-ctx.Done():
			return
		case results <- result:
		}
	}
}

// NewCommand initializes and returns a exec.Cmd for the command provided.
func NewCommand(ctx context.Context, stdout, stderr io.Writer, workingDir string, args ...string) (*exec.Cmd, error) {
	var cmdArgs []string
	if runtime.GOOS == "windows" {
		cmdArgs = []string{"cmd.exe", "/C"}
	} else {
		cmdArgs = []string{"sh", "-c"}
	}

	if len(args) == 0 {
		return nil, fmt.Errorf("failed to create command, no arguments provided")
	}

	cmdArgs = append(cmdArgs, args...)
	cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = os.Environ()
	cmd.Dir = workingDir

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd, nil
}
