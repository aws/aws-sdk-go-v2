package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/changelog"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
)

const listHelpDoc = `changelog ls
`

var listFlagSet = func() *flag.FlagSet {
	fs := flag.NewFlagSet("ls", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprint(fs.Output(), listHelpDoc)
	}
	return fs
}()

func runListCommand(args []string, repoRoot string) error {
	if err := listFlagSet.Parse(args); err != nil {
		return err
	}

	annotations, err := changelog.GetAnnotations(repoRoot)
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"ID", "Type", "Modules", "Collapse", "Description"})

	for _, annotation := range annotations {
		table.Append([]string{annotation.ID, annotation.Type.String(), strconv.Itoa(len(annotation.Modules)), strconv.FormatBool(annotation.Collapse), annotation.Description})
	}

	table.Render()

	return nil
}
