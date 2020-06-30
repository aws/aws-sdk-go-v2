package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/internal/tools/changes"
	"strconv"
)

var changeParams = struct {
	module      string
	changeType  string
	description string
	similar     bool
}{}

var addFlags *flag.FlagSet
var lsFlags *flag.FlagSet

func init() {
	addFlags = flag.NewFlagSet("add", flag.ExitOnError)
	addFlags.StringVar(&changeParams.module, "module", "", "creates a change for the specified module")
	addFlags.StringVar(&changeParams.changeType, "type", "", "sets the change's type")
	addFlags.StringVar(&changeParams.description, "description", "", "sets the change's description")

	lsFlags = flag.NewFlagSet("ls", flag.ExitOnError)
	lsFlags.StringVar(&changeParams.module, "module", "", "filters changes by module")
}

func changeSubcmd(args []string) error {
	if len(args) == 0 {
		usage()
	}

	changesPath, err := changes.GetChangesPath()
	if err != nil {
		return fmt.Errorf("failed to load .changes directory: %v", err)
	}

	metadata, err := changes.LoadMetadata(changesPath)
	if err != nil {
		return fmt.Errorf("failed to load .changes directory: %v", err)
	}

	switch args[0] {
	case "add", "new":
		addFlags.Parse(args[1:])
		return addCmd(metadata, changeParams.module, changeParams.changeType, changeParams.description)
	case "ls", "list":
		lsFlags.Parse(args[1:])
		return lsCmd(metadata, changeParams.module)
	case "modify", "edit":
		if len(args) < 2 {
			usage()
		}

		return modifyCmd(metadata, args[1])
	case "rm", "delete":
		if len(args) < 2 {
			usage()
		}

		return rmCmd(metadata, args[1])
	default:
		usage()
	}

	return nil
}

func addCmd(metadata *changes.Metadata, module, changeType, description string) error {
	if module == "" {
		currentModule, err := changes.GetCurrentModule()
		if err != nil {
			return fmt.Errorf("failed to create change: the module flag was not provided and the tool could not detect a module")
		}

		module = currentModule
	}

	var newChanges []*changes.Change
	var err error

	if changeType != "" && description != "" {
		newChanges, err = changes.NewChanges([]string{module}, changeType, description)
		if err != nil {
			return fmt.Errorf("failed to create change: %v", err)
		}

		err = metadata.AddChanges(newChanges)
		if err != nil {
			return fmt.Errorf("failed to create change: %v", err)
		}
	} else {
		template, err := changes.ChangeToTemplate(&changes.Change{
			Module: module,
		})
		if err != nil {
			return fmt.Errorf("failed to create change: %v", err)
		}

		filledTemplate, err := editTemplate(template)
		if err != nil {
			return fmt.Errorf("failed to create change: %v", err)
		}

		newChanges, err = metadata.AddChangesFromTemplate(filledTemplate)
		if err != nil {
			return fmt.Errorf("failed to create change: %v", err)
		}
	}

	for _, c := range newChanges {
		fmt.Println("added change with id " + c.ID)
	}

	return nil
}

func lsCmd(metadata *changes.Metadata, module string) error {
	for i, c := range metadata.Changes {
		if c.Module == module || module == "" {
			fmt.Printf("[%d] %s\n", i, c.ID)
			fmt.Println("\t", c.Type)
			fmt.Println("\t", c.Description)
			fmt.Println()
		}
	}

	return nil
}

func modifyCmd(metadata *changes.Metadata, id string) error {
	change, err := selectChange(metadata, id)
	if err != nil {
		return fmt.Errorf("failed to modify change: %v", err)
	}

	template, err := changes.ChangeToTemplate(change)
	if err != nil {
		return fmt.Errorf("failed to modify change: %v", err)
	}

	filledTemplate, err := editTemplate(template)
	if err != nil {
		return fmt.Errorf("failed to modify change: %v", err)
	}

	newChange, err := metadata.UpdateChangeFromTemplate(change, filledTemplate)
	if err != nil {
		return fmt.Errorf("couldn't modify change: %v", err)
	}

	fmt.Printf("successfully modified %s, new id is %s\n", change.ID, newChange.ID)
	return nil
}

func rmCmd(metadata *changes.Metadata, id string) error {
	change, err := selectChange(metadata, id)
	if err != nil {
		return fmt.Errorf("failed to remove change: %v", err)
	}

	err = metadata.RemoveChangeById(change.ID)
	if err != nil {
		return fmt.Errorf("failed to remove change: %v", err)
	}

	fmt.Println("successfully removed " + change.ID)
	return nil
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
