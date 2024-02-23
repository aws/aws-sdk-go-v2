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
	maxDepth int = 6
	serviceDir string
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

	currentDepth := strings.Count(path, string(os.PathSeparator))

	isInternal := strings.Count(path, "/internal") > 0
	if isInternal {
		fmt.Printf("Skipping %v\n", path)
		return nil
	}

	if currentDepth > maxDepth || isInternal {
		return fs.SkipDir
	}

	items := map[string]JewelryItem{}
	Extract(path, d, items)
	fmt.Printf("processed %v\n", path)
	return nil
}

func init() {
	flag.StringVar(&serviceDir, "serviceDir", "",
		"Root directory that is direct parent of all service client directories")

}

func main() {
	flag.Parse()
	args := flag.Args()
	serviceDir := args[0]
	fmt.Println(serviceDir)
	if serviceDir == "" {
		log.Fatalf("no service directory specified")
	}

	err := filepath.WalkDir(serviceDir, extract)
	if err != nil {
		log.Fatal(err)
	}
}
