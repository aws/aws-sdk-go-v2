package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/changelog"
	"os"
)

const viewHelpDoc = `changelog view <id>
`

var viewFlagSet = func() *flag.FlagSet {
	fs := flag.NewFlagSet("view", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprint(fs.Output(), viewHelpDoc)
	}
	return fs
}()

func runViewCommand(args []string, repoRoot string) error {
	if err := viewFlagSet.Parse(args); err != nil {
		return err
	}

	args = viewFlagSet.Args()

	if len(args) == 0 || len(args) > 1 {
		return fmt.Errorf("expect one annotation id to be provided")
	}

	annotation, err := changelog.LoadAnnotation(repoRoot, args[0])
	if err != nil {
		return err
	}

	marshal, err := json.MarshalIndent(annotation, "", "    ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(os.Stdout, "%s\n", marshal)
	return err
}
