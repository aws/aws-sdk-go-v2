package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/internal/repotools/changes"
)

var changeParams = struct {
	module      string
	changeType  changes.ChangeType
	description string
	compareTo   string
	similar     bool
}{}

var addFlags *flag.FlagSet
var lsFlags *flag.FlagSet
var modifyFlags *flag.FlagSet
var rmFlags *flag.FlagSet

func changeUsage() {
	sets := []*flag.FlagSet{addFlags, lsFlags, modifyFlags, rmFlags}

	for _, f := range sets {
		f.Usage()
	}
}

func init() {
	addFlags = flag.NewFlagSet("add", flag.ExitOnError)
	addFlags.StringVar(&changeParams.module, "module", "", "sets the change's module")
	addFlags.Var(&changeParams.changeType, "type", "sets the change's type")
	addFlags.StringVar(&changeParams.description, "description", "", "sets the change's description")
	addFlags.StringVar(&changeParams.compareTo, "compare-to", "", "specifies a path to a version enclosure to compare current module hashes to in order to resolve a wildcard.")
	addFlags.Usage = func() {
		fmt.Printf("%s change add [-module=<module>] [-type=<type>] [-description=<description>]\n", os.Args[0])
		addFlags.PrintDefaults()
	}

	lsFlags = flag.NewFlagSet("ls", flag.ExitOnError)
	lsFlags.StringVar(&changeParams.module, "module", "", "filters changes by module")
	lsFlags.Usage = func() {
		fmt.Printf("%s change ls [-module=<module>]\n", os.Args[0])
		lsFlags.PrintDefaults()
	}

	modifyFlags = flag.NewFlagSet("modify", flag.ExitOnError)
	modifyFlags.Usage = func() {
		fmt.Printf("%s change modify <change id>\n  <change id>: the index (as found in the ls subcommand) or the ID of the change to modify\n", os.Args[0])
		modifyFlags.PrintDefaults()
	}

	rmFlags = flag.NewFlagSet("rm", flag.ExitOnError)
	rmFlags.Usage = func() {
		fmt.Printf("%s change rm <change id>\n  <change id>: the index (as found in the ls subcommand) or the ID of the change to remove\n", os.Args[0])
		rmFlags.PrintDefaults()
	}
}

func changeSubcmd(args []string) error {
	if len(args) == 0 {
		changeUsage()
		return errors.New("invalid usage")
	}

	subCommand := args[0]

	changesPath, err := changes.GetChangesPath()
	if err != nil {
		return fmt.Errorf("failed to load .changes directory: %v", err)
	}

	metadata, err := changes.LoadMetadata(changesPath)
	if err != nil {
		return fmt.Errorf("failed to load .changes directory: %v", err)
	}

	switch subCommand {
	case "add", "new":
		err = addFlags.Parse(args[1:])
		if err != nil {
			return err
		}

		if changes.ModIsWildcard(changeParams.module) {
			return addCmdWildcard(metadata, changeParams.module, changeParams.changeType, changeParams.description, changeParams.compareTo)
		}

		return addCmd(metadata, changeParams.module, changeParams.changeType, changeParams.description)
	case "ls", "list":
		err = lsFlags.Parse(args[1:])
		if err != nil {
			return err
		}

		return lsCmd(metadata, changeParams.module)
	case "modify", "edit":
		err = modifyFlags.Parse(args[1:])
		if err != nil {
			return err
		}

		if len(args) < 2 {
			changeUsage()
			return errors.New("invalid usage")
		}

		id := args[1]

		return modifyCmd(metadata, id)
	case "rm", "delete":
		err = rmFlags.Parse(args[1:])
		if err != nil {
			return err
		}

		if len(args) < 2 {
			changeUsage()
			return errors.New("invalid usage")
		}

		id := args[1]

		return rmCmd(metadata, id)
	default:
		changeUsage()
		return errors.New("invalid usage")
	}
}

func addCmdWildcard(metadata *changes.Metadata, module string, changeType changes.ChangeType, description string, compareTo string) error {
	if module == "" {
		return errors.New("couldn't add wildcard change: a module must be provided with --module")
	}

	repo, err := changes.NewRepository(filepath.Join(metadata.ChangePath, ".."))
	if err != nil {
		return fmt.Errorf("couldn't add wildcard change: %v", err)
	}

	mods, err := repo.Modules()
	if err != nil {
		return err
	}

	var affectedModules []string
	if compareTo != "" {
		data, err := ioutil.ReadFile(compareTo)
		if err != nil {
			return err
		}

		var enc changes.VersionEnclosure
		err = json.Unmarshal(data, &enc)
		if err != nil {
			return err
		}

		hashes, err := repo.ModuleHashes(enc)
		if err != nil {
			return err
		}

		affectedModules = enc.HashDiff(hashes)
	} else {
		affectedModules, err = changes.MatchWildcardModules(mods, module)
		if err != nil {
			return err
		}
	}

	template, err := changes.ChangeToTemplate(changes.Change{
		Module:          module,
		Type:            changeType,
		Description:     description,
		AffectedModules: affectedModules,
	})
	if err != nil {
		return fmt.Errorf("failed to create change: %v", err)
	}

	filledTemplate, err := editTemplate(template)
	if err != nil {
		return fmt.Errorf("failed to create change: %v", err)
	}

	changes, err := changes.TemplateToChanges(filledTemplate)
	if err != nil {
		return fmt.Errorf("failed to create change: %v", err)
	}

	if len(changes) != 1 {
		return fmt.Errorf("failed to create change: expected template to create 1 change, got %d changes", len(changes))
	}

	change := changes[0]

	// TODO: move some logic into TemplateToChanges

	return metadata.AddChange(change)
}

func addCmd(metadata *changes.Metadata, module string, changeType changes.ChangeType, description string) error {
	if module == "" {
		currentModule, err := changes.GetCurrentModule()
		if err != nil {
			return fmt.Errorf("failed to create change: the module flag was not provided and the tool could not detect a module")
		}

		module = currentModule
	}

	var newChanges []changes.Change
	var err error

	if changeType != "" && description != "" {
		newChanges, err = changes.NewChanges([]string{module}, changeType, description, "")
		if err != nil {
			return fmt.Errorf("failed to create change: %v", err)
		}

		err = metadata.AddChanges(newChanges)
		if err != nil {
			return fmt.Errorf("failed to create change: %v", err)
		}
	} else {
		template, err := changes.ChangeToTemplate(changes.Change{
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

	newChanges, err := metadata.UpdateChangeFromTemplate(change, filledTemplate)
	if err != nil {
		return fmt.Errorf("couldn't modify change: %v", err)
	}

	fmt.Printf("successfully modified %s, new change(s):\n", change.ID)
	for _, c := range newChanges {
		fmt.Printf("\t%s\n", c.ID)
	}
	return nil
}

func rmCmd(metadata *changes.Metadata, id string) error {
	change, err := selectChange(metadata, id)
	if err != nil {
		return fmt.Errorf("failed to remove change: %v", err)
	}

	err = metadata.RemoveChangeByID(change.ID)
	if err != nil {
		return fmt.Errorf("failed to remove change: %v", err)
	}

	fmt.Println("successfully removed " + change.ID)
	return nil
}

// selectChange will return the change identified by the given id, which can be either the index of one of metadata's
// Changes or the Change's ID.
func selectChange(metadata *changes.Metadata, id string) (changes.Change, error) {
	// try selecting by index first
	index, err := strconv.Atoi(id)
	if err == nil {
		if index < 0 || index >= len(metadata.Changes) {
			return changes.Change{}, fmt.Errorf("failed to get change with index %d: index out of range", index)
		}
		return metadata.Changes[index], nil
	}

	return metadata.GetChangeByID(id)
}
