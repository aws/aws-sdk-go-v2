package main

import (
	"flag"
	"fmt"
	"github.com/aggagen/changes"
	"log"
	"strconv"
)

var changeParams = struct {
	module  string
	similar bool
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
	}

	changesPath, err := changes.GetChangesPath()
	if err != nil {
		log.Fatalf("Failed to load .changes directory: %v", err)
	}

	metadata, err := changes.LoadMetadata(changesPath)
	if err != nil {
		log.Fatalf("Failed to load .changes directory: %v", err)
	}

	switch args[0] {
	case "add", "new":
		addFlags.Parse(args[1:])
		addCmd(metadata, changeParams.module)
	case "ls", "list":
		lsFlags.Parse(args[1:])
		lsCmd(metadata, changeParams.module)
	case "modify", "edit":
		if len(args) < 2 {
			usage()
		}

		modifyCmd(metadata, args[1])
	case "rm", "delete":
		if len(args) < 2 {
			usage()
		}

		rmCmd(metadata, args[1])
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
		log.Fatalf("failed to create change: %v", err)
	}

	filledTemplate, err := editTemplate(template)
	if err != nil {
		log.Fatalf("failed to create change: %v", err)
	}

	newChanges, err := metadata.AddChangesFromTemplate(filledTemplate)
	if err != nil {
		log.Fatalf("failed to create change: %v", err)
	}

	for _, c := range newChanges {
		fmt.Println("added change with id " + c.ID)
	}
}

func lsCmd(metadata *changes.Metadata, module string) {
	for i, c := range metadata.Changes {
		if c.Module == module || module == "" {
			fmt.Printf("[%d] %s\n", i, c.ID)
			fmt.Println("\t", c.Type)
			fmt.Println("\t", c.Description)
			fmt.Println()
		}
	}
}

func modifyCmd(metadata *changes.Metadata, id string) {
	change, err := selectChange(metadata, id)
	if err != nil {
		log.Fatalf("failed to modify change: %v", err)
	}

	template, err := changes.ChangeToTemplate(change)
	if err != nil {
		log.Fatalf("failed to modify change: %v", err)
	}

	filledTemplate, err := editTemplate(template)
	if err != nil {
		log.Fatalf("failed to modify change: %v", err)
	}

	newChange, err := metadata.UpdateChangeFromTemplate(change, filledTemplate)
	if err != nil {
		log.Fatalf("couldn't modify change: %v", err)
	}

	fmt.Printf("successfully modified %s, new id is %s\n", change.ID, newChange.ID)
}

func rmCmd(metadata *changes.Metadata, id string) {
	change, err := selectChange(metadata, id)
	if err != nil {
		log.Fatalf("failed to remove change: %v", err)
	}

	err = metadata.RemoveChangeById(change.ID)
	if err != nil {
		log.Fatalf("failed to remove change: %v", err)
	}

	fmt.Println("successfully removed " + change.ID)
}

// selectChange will return the change identified by the given id, which can be either the index of one of metadata's
// Changes or the Change's ID.
func selectChange(metadata *changes.Metadata, id string) (*changes.Change, error) {
	// try selecting by index first
	index, err := strconv.Atoi(id)
	if err == nil {
		if index < 0 || index >= len(metadata.Changes) {
			return nil, fmt.Errorf("failed to get change with index %d: index out of range\n", index)
		}
		return metadata.Changes[index], nil
	}

	return metadata.GetChangeById(id)
}
