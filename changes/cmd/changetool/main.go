package main

import (
	"flag"
	"fmt"
	"github.com/aggagen/changes"
	"os"
)

func usage() {
	fmt.Printf("USAGE: todo")
}

func main() {
	metadata, err := changes.LoadMetadata()
	if err != nil {
		fmt.Printf("Failed to load .changes directory: %v\n", err)
		os.Exit(1)
	}

	addFlags := flag.NewFlagSet("add", flag.ExitOnError)
	addModule := addFlags.String("module", "", "creates a change for the specified module")

	lsFlags := flag.NewFlagSet("ls", flag.ExitOnError)
	lsModule := lsFlags.String("module", "", "filters changes by module")

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addFlags.Parse(os.Args[2:])
		addCmd(metadata, *addModule)
	case "ls", "list":
		lsFlags.Parse(os.Args[2:])
		lsCmd(metadata, *lsModule)
	case "modify", "edit":
		if len(os.Args) < 3 {
			usage()
			os.Exit(1)
		}

		modifyCmd(metadata, os.Args[2])
	case "rm", "delete":
		if len(os.Args) < 3 {
			usage()
			os.Exit(1)
		}

		rmCmd(metadata, os.Args[2])
	case "release":
		fmt.Println(metadata.CreateRelease("2020-06-19", []changes.VersionBump{}))
	}
}

func addCmd(metadata *changes.Metadata, module string) {
	if module == "" {
		if currentModule, err := changes.GetCurrentModule(); err == nil {
			module = currentModule
		}
	}

	filledTemplate, err := editTemplate(changes.ChangeToTemplate(&changes.Change{
		Modules: []string{module},
	}))
	if err != nil {
		fmt.Printf("failed to create change: %v\n", err)
		os.Exit(1)
	}

	change := &changes.Change{}

	err = changes.TemplateToChange(filledTemplate, change)
	if err != nil {
		fmt.Printf("failed to create change: %v\n", err)
		os.Exit(1)
	}

	err = metadata.AddChange(change)
	if err != nil {
		fmt.Printf("failed to create change: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("added change with id " + change.Id)
}

func lsCmd(metadata *changes.Metadata, module string) {
	for _, c := range metadata.Changes {
		if module != "" && !c.AffectsModule(module) {
			continue
		}

		fmt.Println(c.Id)
		fmt.Println("\t", c.Type)
		fmt.Println("\t", c.Description)
		fmt.Println()
	}
}

func modifyCmd(metadata *changes.Metadata, id string) {
	change, err := metadata.GetChangeById(id)
	if err != nil {
		fmt.Printf("failed to modify change: %v\n", err)
		os.Exit(1)
	}


	filledTemplate, err := editTemplate(changes.ChangeToTemplate(change))
	if err != nil {
		fmt.Printf("failed to modify change: %v\n", err)
		os.Exit(1)
	}

	err = changes.TemplateToChange(filledTemplate, change)
	if err != nil {
		fmt.Printf("failed to modify change: %v\n", err)
		os.Exit(1)
	}

	err = metadata.SaveChange(change)
	if err != nil {
		fmt.Printf("failed to modify change: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("successfully modified " + id)
}

func rmCmd(metadata *changes.Metadata, id string) {
	err := metadata.RemoveChangeById(id)
	if err != nil {
		fmt.Printf("failed to remove change: %v", err)
		os.Exit(1)
	}

	fmt.Println("successfully removed " + id)
}