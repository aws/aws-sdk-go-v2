package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"

	repotools "github.com/awslabs/aws-go-multi-module-repository-tools"
)

var (
	atOnce             int
	rootPath           string
	pathRelRoot        bool
	skipRootPath       bool
	skipPaths          string
	skipEmptyRootPaths bool
	failFast           bool
)

func init() {
	flag.BoolVar(&failFast, "fail-fast", true,
		"Terminates the module walking and command as soon as a error in any command occurs.")

	flag.BoolVar(&skipEmptyRootPaths, "skip-empty-root", true,
		"Directs to skip the root path if empty.")

	flag.BoolVar(&skipRootPath, "skip-root", false,
		"Directs to skip the `root path` and only run commands on discovered submodules.")

	flag.BoolVar(&pathRelRoot, "rel-repo", true,
		"Directs if the path is relative to the repository (true) or working directory (false).")

	flag.StringVar(&rootPath, "p", "",
		"The root `path` to walk each module from. If unset walks to the repository root.")

	flag.IntVar(&atOnce, "c", 1,
		"Number of `concurrent commands` to invoke at once.")

	flag.StringVar(&skipPaths, "skip", "",
		"Set of `paths to skip`, delimited with "+string(os.PathListSeparator))
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

		} else if !filepath.IsAbs(rootPath) {
			rootPath, err = repotools.JoinWorkingDirectory(rootPath)
			if err != nil {
				return fmt.Errorf("failed to get relative path, %w", err)
			}
		}

	}

	// Skip built in paths relative from the repo root.
	for _, skip := range getSkipDirs() {
		boots.SkipDirs = append(boots.SkipDirs, filepath.Join(repoRoot, skip))
	}

	// Skip additional paths relative to the root path.
	for _, skip := range strings.Split(skipPaths, string(os.PathListSeparator)) {
		skip = strings.TrimSpace(skip)
		if len(skip) == 0 {
			continue
		}
		boots.SkipDirs = append(boots.SkipDirs, filepath.Join(rootPath, skip))
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

	// Logging command status
	var failed bool
	var resWG sync.WaitGroup
	resWG.Add(1)
	results := make(chan WorkLog)
	go func() {
		defer resWG.Done()
		for result := range results {
			relPath, err := filepath.Rel(repoRoot, result.Path)
			if err != nil {
				log.Println("failed to get path relative to repo root",
					repoRoot, result.Path, err)
				relPath = result.Path
			}

			var output string
			if result.Output != nil {
				b, err := ioutil.ReadAll(result.Output)
				if err != nil {
					log.Printf("%s: %s => failed to read result output, %v",
						relPath, result.Cmd, err)
				}
				output = string(b)
			}

			if result.Err != nil {
				log.Printf("%s: %s => error: %v\n%s",
					relPath, result.Cmd, result.Err, output)
				failed = true
			} else {
				log.Printf("%s: %s =>\n%s",
					relPath, result.Cmd, output)
			}

			//  Terminate early as soon as any command fails.
			if failFast && result.Err != nil {
				cancelFn()
			}
		}
	}()

	// Work consumer
	var jobWG sync.WaitGroup
	jobWG.Add(atOnce)
	jobs := make(chan Work)
	for i := 0; i < atOnce; i++ {
		go func() {
			defer jobWG.Done()
			var streamOut io.Writer
			if atOnce == 1 {
				streamOut = os.Stdout
			}
			CommandWorker(ctx, jobs, results, streamOut)
		}()
	}

	// Special case to skip root path when path if they don't contain go files.
	if skipEmptyRootPaths {
		matches, err := filepath.Glob(filepath.Join(rootPath, "*.go"))
		if err != nil || len(matches) == 0 {
			skipRootPath = true
		}
	}

	modulePaths := boots.Modules()
	if skipRootPath {
		modulePaths = removePath(rootPath, modulePaths)

	} else if !hasPath(rootPath, modulePaths) {
		modulePaths = append([]string{rootPath}, modulePaths...)
	}

	// Work producer
Loop:
	for _, modPath := range modulePaths {
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

func hasPath(path string, paths []string) bool {
	for i := 0; i < len(paths); i++ {
		if paths[i] == path {
			return true
		}
	}
	return false
}

func removePath(path string, paths []string) []string {
	for i := 0; i < len(paths); i++ {
		if paths[i] == path {
			paths = append(paths[:i], paths[i+1:]...)
			i--
		}
	}
	return paths
}
