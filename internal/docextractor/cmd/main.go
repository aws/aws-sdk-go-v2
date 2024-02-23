package main

import (
	"fmt"
	"flag"
	"os"
	"log"
	"path/filepath"
	"strings"
	"io/fs"
)

var (
	maxDepth int
	servicePath string
)

// Visit function passed to filepath.WalkDir. 
// Ensures that Exctract is only called on service client
// directories and not internal directories or files.
func extract(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !d.IsDir() {
		return nil
	}

	isInternal := strings.Count(path, "/internal") > 0
	if isInternal {
		fmt.Printf("Skipping %v\n", path)
		return nil
	}

	currentDepth := strings.Count(path, string(os.PathSeparator))
	serviceDepth := strings.Count(servicePath, string(os.PathSeparator))

	if currentDepth > (serviceDepth + 1) || isInternal {
		return fs.SkipDir
	}

	items := map[string]JewelryItem{}
	Extract(path, d, items)
	fmt.Printf("processed %v\n", path)
	return nil
}

func init() {
	flag.StringVar(&servicePath, "servicePath", "",
		"Root directory that is direct parent of all service client directories")

}

func main() {
	flag.Parse()
	if servicePath == "" {
		log.Fatalf("need service dir and max depth")
	}

	log.Println(
		fmt.Sprintf("Processing service path %v", servicePath),
	)


	err := filepath.WalkDir(servicePath, extract)
	if err != nil {
		log.Fatal(err)
	}
}
