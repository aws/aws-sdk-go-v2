package main

import (
	"flag"
	"fmt"
	"github.com/aggagen/changes"
	"os"
	"strconv"
)

var changeParams = struct {
	module string
}{}

var addFlags *flag.FlagSet
var lsFlags *flag.FlagSet

func init() {
	addFlags = flag.NewFlagSet("add", flag.ExitOnError)
	addFlags.StringVar(&changeParams.module, "module", "", "creates a change for the specified module")

	lsFlags = flag.NewFlagSet("ls", flag.ExitOnError)
	lsFlags.StringVar(&changeParams.module, "module", "", "filters changes by module")
}

func changeSubcmd(args []string) {
	if len(args) == 0 {
		usage()
		os.Exit(1)
	}

	changesPath, err := changes.GetChangesPath()
	if err != nil {
		fmt.Printf("Failed to load .changes directory: %v\n", err)
		os.Exit(1)
	}

	metadata, err := changes.LoadMetadata(changesPath)
	if err != nil {
		fmt.Printf("Failed to load .changes directory: %v\n", err)
		os.Exit(1)
	}

	switch args[0] {
	case "add":
		addFlags.Parse(args[1:])
		addCmd(metadata, changeParams.module)
	case "ls", "list":
		lsFlags.Parse(args[1:])
		lsCmd(metadata, changeParams.module)
	case "modify", "edit":
		if len(args) < 2 {
			usage()
			os.Exit(1)
		}

		modifyCmd(metadata, args[1])
	case "rm", "delete":
		if len(args) < 2 {
			usage()
			os.Exit(1)
		}

		rmCmd(metadata, os.Args[1])
	}
}

func addCmd(metadata *changes.Metadata, module string) {
	if module == "" {
		currentModule, err := changes.GetCurrentModule()
		if err != nil {
			fmt.Printf("failed to create change: the module flag was not provided and the tool could not detect a module")
		}

		module = currentModule
	}

	template, err := changes.ChangeToTemplate(&changes.Change{
		Module: module,
	})
	if err != nil {
		fmt.Printf("failed to create change: %v\n", err)
		os.Exit(1)
	}

	filledTemplate, err := editTemplate(template)
	if err != nil {
		fmt.Printf("failed to create change: %v\n", err)
		os.Exit(1)
	}

	newChanges, err := changes.TemplateToChanges(filledTemplate)
	if err != nil {
		fmt.Printf("failed to create change: %v\n", err)
		os.Exit(1)
	}

	err = metadata.AddChanges(newChanges)
	if err != nil {
		fmt.Printf("failed to create change: %v\n", err)
		os.Exit(1)
	}

	for _, c := range newChanges {
		fmt.Println("added change with id " + c.ID)
	}
}

func lsCmd(metadata *changes.Metadata, module string) {
	for i, c := range metadata.ListChanges() {
		if module != "" && module != c.Module {
			continue
		}

		fmt.Printf("[%d] %s\n", i, c.ID)
		fmt.Println("\t", c.Type)
		fmt.Println("\t", c.Description)
		fmt.Println()
	}
}

func modifyCmd(metadata *changes.Metadata, id string) {
	var change *changes.Change
	index, err := strconv.Atoi(id)
	if err == nil {
		if index < 0 || index >= len(metadata.ListChanges()) {
			fmt.Printf("failed to modify change with index %d: index out of range\n", index)
			os.Exit(1)
		}
		change = metadata.ListChanges()[index]
	} else {
		change, err = metadata.GetChangeById(id)
		if err != nil {
			fmt.Printf("failed to modify change: %v\n", err)
			os.Exit(1)
		}
	}

	template, err := changes.ChangeToTemplate(change)
	if err != nil {
		fmt.Printf("failed to modify change: %v\n", err)
		os.Exit(1)
	}

	filledTemplate, err := editTemplate(template)
	if err != nil {
		fmt.Printf("failed to modify change: %v\n", err)
		os.Exit(1)
	}

	if template == filledTemplate {
		fmt.Println("no change was made to " + change.ID)
		os.Exit(1)
	}

	newChanges, err := changes.TemplateToChanges(filledTemplate)
	if err != nil {
		fmt.Printf("failed to modify change: %v\n", err)
		os.Exit(1)
	} else if len(newChanges) != 1 {
		fmt.Printf("failed to modify change: modules cannot be added to a change during modification")
		os.Exit(1)
	}

	err = metadata.SaveChange(newChanges[0])
	if err != nil {
		fmt.Printf("failed to modify change: %v\n", err)
		os.Exit(1)
	}

	err = metadata.RemoveChangeById(change.ID)
	if err != nil {
		fmt.Printf("failed to remove old change with id %s: %v\n", change.ID, err)
		os.Exit(1)
	}

	fmt.Printf("successfully modified %s, new id is %s\n", change.ID, newChanges[0].ID)
}

func rmCmd(metadata *changes.Metadata, id string) {
	err := metadata.RemoveChangeById(id)
	if err != nil {
		fmt.Printf("failed to remove change: %v", err)
		os.Exit(1)
	}

	fmt.Println("successfully removed " + id)
}
