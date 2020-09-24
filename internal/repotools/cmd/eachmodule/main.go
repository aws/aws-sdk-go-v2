package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go-v2/internal/repotools"
)

var (
	atOnce      int
	rootPath    string
	pathRelRoot bool
)

func init() {
	flag.BoolVar(&pathRelRoot, "rel-repo", true,
		"Directs if the path is relative the repository (true) or working directory (false).")
	flag.StringVar(&rootPath, "p", "",
		"The root `path` to walk each module from. If unset walks to the repository root.")
	flag.IntVar(&atOnce, "c", 1,
		"Number of `concurrent` commands to invoke at once.")

	// TODO add skip dirs relative to root (or repo root if root isn't set)
}

// SkipDir paths are all relative to the root of the repository.
func getSkipDirs() []string {
	return []string{
		"codegen",
	}
}

func run() (err error) {
	flag.Parse()
	cmds := flag.Args()
	if len(cmds) == 0 {
		log.Fatalf("no command specified")
	}

	var boots repotools.Boots

	repoRoot, err := repotools.FindRepoRoot(rootPath)
	if err != nil {
		return fmt.Errorf("failed to get repository root path, %w", err)
	}

	if pathRelRoot {
		rootPath = filepath.Join(repoRoot, rootPath)

	} else {
		if len(rootPath) == 0 {
			rootPath, err = repotools.FindRepoRoot(rootPath)
			if err != nil {
				return fmt.Errorf("failed to get repository root path, %w", err)
			}

			for _, skip := range getSkipDirs() {
				boots.SkipDirs = append(boots.SkipDirs, filepath.Join(rootPath, skip))
			}
		} else if !filepath.IsAbs(rootPath) {
			rootPath, err = repotools.JoinWorkingDirectory(rootPath)
			if err != nil {
				return fmt.Errorf("failed to get relative path, %w", err)
			}
		}
	}

	if err := filepath.Walk(rootPath, boots.Walk); err != nil {
		return fmt.Errorf("failed to walk directory, %w", err)
	}

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	// Block until a signal is received.
	go func() {
		<-c
		cancelFn()
	}()

	var failed bool
	var resWG sync.WaitGroup
	resWG.Add(1)
	results := make(chan WorkLog)
	go func() {
		defer resWG.Done()
		for result := range results {
			if result.Err != nil {
				failed = true
			}
			relPath, err := filepath.Rel(repoRoot, result.Path)
			if err != nil {
				log.Println("failed to get path relative to repo root",
					repoRoot, result.Path, err)
				relPath = result.Path
			}
			log.Printf("%s: %s => error: %v\n%s",
				relPath, result.Cmd, result.Err, result.Output.String())
		}
	}()

	var jobWG sync.WaitGroup
	jobWG.Add(atOnce)
	jobs := make(chan Work)
	for i := 0; i < atOnce; i++ {
		go func() {
			defer jobWG.Done()
			CommandWorker(ctx, jobs, results)
		}()
	}

Loop:
	for _, modPath := range boots.Modules() {
		for _, cmd := range cmds {
			select {
			case <-ctx.Done():
				break Loop
			case jobs <- Work{
				Path: modPath,
				Cmd:  cmd,
			}:
			}
		}
	}
	close(jobs)

	jobWG.Wait()
	close(results)

	resWG.Wait()

	if failed {
		return fmt.Errorf("a command failed")
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
