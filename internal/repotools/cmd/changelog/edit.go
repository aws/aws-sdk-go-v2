package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/changelog"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/gomod"
)

const editHelpDoc = `changelog edit <id>
`

var editFlagSet = func() *flag.FlagSet {
	fs := flag.NewFlagSet("edit", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprint(fs.Output(), editHelpDoc)
	}
	return fs
}()

func runEditCommand(args []string, repoRoot string) error {
	if err := editFlagSet.Parse(args); err != nil {
		return err
	}

	discoverer := gomod.NewDiscoverer(repoRoot)

	if err := discoverer.Discover(); err != nil {
		return err
	}

	modules, err := discoverer.ModulesRel()
	if err != nil {
		return fmt.Errorf("failed to discover repository go modules: %w", err)
	}

	args = editFlagSet.Args()
	if len(args) == 0 || len(args) > 1 {
		return fmt.Errorf("expect one annotation id to be provided")
	}

	annotation, err := changelog.LoadAnnotation(repoRoot, args[0])
	if err != nil {
		return err
	}

	if err := interactiveEdit(&annotation, modules); err != nil {
		return err
	}

	return changelog.WriteAnnotation(repoRoot, annotation)
}
