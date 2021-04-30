package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/internal/repotools"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/changelog"
)

// ChangeType is an alias type for changelog.ChangeType to add flag support
type ChangeType changelog.ChangeType

func (c *ChangeType) String() string {
	return changelog.ChangeType(*c).String()
}

// Set parses c and sets it to the changelog.ChangeType
func (c *ChangeType) Set(s string) error {
	ct := changelog.ParseChangeType(s)
	*c = ChangeType(ct)
	return nil
}

func init() {
	flag.Usage = printHelp
}

func main() {
	flag.Parse()

	repoRoot, err := repotools.GetRepoRoot()
	if err != nil {
		log.Fatalf("failed to get repository root: %v", err)
	}

	arg := flag.Arg(0)

	switch {
	case strings.EqualFold(arg, createFlagSet.Name()):
		err = runCreateCommand(flag.Args()[1:], repoRoot)
	case strings.EqualFold(arg, listFlagSet.Name()):
		err = runListCommand(flag.Args()[1:], repoRoot)
	case strings.EqualFold(arg, viewFlagSet.Name()):
		err = runViewCommand(flag.Args()[1:], repoRoot)
	case strings.EqualFold(arg, editFlagSet.Name()):
		err = runEditCommand(flag.Args()[1:], repoRoot)
	case strings.EqualFold(arg, removeFlagSet.Name()):
		err = runRemoveCommand(flag.Args()[1:], repoRoot)
	case strings.EqualFold(arg, "help") || len(arg) == 0:
		fallthrough
	default:
		printHelp()
		return
	}

	if err != nil {
		log.Fatal(err)
	}
}

func printHelp() {
	var builder strings.Builder
	builder.WriteString("Usage:\n\n")
	builder.WriteString(createHelpDoc)
	builder.WriteRune('\n')
	builder.WriteString(listHelpDoc)
	builder.WriteRune('\n')
	builder.WriteString(editHelpDoc)
	builder.WriteRune('\n')
	builder.WriteString(viewHelpDoc)
	builder.WriteRune('\n')
	fmt.Fprint(os.Stderr, builder.String())
	os.Exit(0)
}
