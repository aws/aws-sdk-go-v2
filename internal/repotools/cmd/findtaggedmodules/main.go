/*
findtaggedmodules finds modules that contain given go build tags.

Given a directory and a list of tags, the command looks at all .go files
inside the directory to find any that contain the tags. Once found, it
finds which module they belong to and returns all directories where any
of the tags are found.

It always returns paths Unix-style with a forward slash ("/") as separator

Usage:

    findtagmodules -tags [tag1,tag2] [-p path]

The flags are:

    -tags
        List of tags to look for, passed as "tag1,tag2"
    -p
        Root path to search from. Default to the repo root
*/
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	repotools "github.com/awslabs/aws-go-multi-module-repository-tools"
)

var (
	tags     string
	rootPath string
)

func init() {
	flag.StringVar(&tags, "tags", "", "Comma-separated list of build tags to search for")
	flag.StringVar(&rootPath, "p", "", "Root path to search from (defaults to repo root)")
}

func main() {
	flag.Parse()
	
	if tags == "" {
		log.Fatal("must specify -tags")
	}

	targetTags := strings.Split(tags, ",")
	for i := range targetTags {
		targetTags[i] = strings.TrimSpace(targetTags[i])
	}

	repoRoot, err := repotools.FindRepoRoot(rootPath)
	if err != nil {
		log.Fatalf("failed to find repo root: %v", err)
	}

	if rootPath == "" {
		rootPath = repoRoot
	} else {
		rootPath = filepath.Join(repoRoot, rootPath)
	}

	var boots repotools.Boots

	if err := filepath.Walk(rootPath, boots.Walk); err != nil {
		log.Fatalf("failed to walk directory: %v", err)
	}

	for _, modPath := range boots.Modules() { 
		if modPath == rootPath {
			continue
		}
		found, err := hasAnyTag(modPath, targetTags)
		if err != nil {
			log.Fatalf("found an error searching for tags: %v", err)
		}
		if found {
			relPath, err := filepath.Rel(repoRoot, modPath)
			// Use Unix style path
			relPath = filepath.ToSlash(relPath)
			if err != nil {
				fmt.Println(modPath)
			} else {
				fmt.Println(relPath)
			}
		}
	}
}

func hasAnyTag(modPath string, targetTags []string) (bool, error) {
	found := false
	err := filepath.WalkDir(modPath, func(path string, info fs.DirEntry, err error) error {
		if err != nil || !strings.HasSuffix(path, ".go") {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Read first two lines only
		buf := make([]byte, 200)
		n, _ := file.Read(buf)
		lines := strings.Split(string(buf[:n]), "\n")
		
		for i := 0; i < 2 && i < len(lines); i++ {
			line := strings.TrimSpace(lines[i])
			if strings.HasPrefix(line, "//go:build") || strings.HasPrefix(line, "// +build") {
				for _, tag := range targetTags {
					if strings.Contains(line, tag) {
						found = true
						return nil
					}
				}
			}
		}
		return nil
	}) 
	return found, err
}
