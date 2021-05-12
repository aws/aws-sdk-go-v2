package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/internal/repotools/changelog"
)

const removeHelpDoc = `changelog rm (-all | <id>...)
`

var removeCommand = struct {
	All bool
}{}

var removeFlagSet = func() *flag.FlagSet {
	fs := flag.NewFlagSet("rm", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprint(fs.Output(), editHelpDoc)
	}
	fs.BoolVar(&removeCommand.All, "all", false, "")
	return fs
}()

func runRemoveCommand(args []string, repoRoot string) error {
	if err := removeFlagSet.Parse(args); err != nil {
		return err
	}

	args = removeFlagSet.Args()
	if len(args) == 0 && !removeCommand.All {
		return fmt.Errorf("expect either a list ids or -all")
	}

	toRemove := make(map[string]struct{})
	for _, id := range args {
		toRemove[id] = struct{}{}
	}

	annotations, err := changelog.GetAnnotations(repoRoot)
	if err != nil {
		return err
	}

	for _, annotation := range annotations {
		_, ok := toRemove[annotation.ID]
		if !removeCommand.All && !ok {
			continue
		}
		if err := changelog.RemoveAnnotation(repoRoot, annotation); err != nil {
			log.Fatalf("failed to remove annotation: %v", err)
		}
	}

	return nil
}
