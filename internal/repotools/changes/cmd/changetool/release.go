package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/internal/repotools/changes"
)

var releaseParams = struct {
	repo           string
	releaseID      string
	selector       changes.VersionSelector
	pretty         bool
	push           bool
	nonInteractive bool
}{}

var staticVersionsFlags *flag.FlagSet
var createReleaseFlags *flag.FlagSet

func releaseUsage() {
	var sets = []*flag.FlagSet{staticVersionsFlags, createReleaseFlags}

	for _, f := range sets {
		f.Usage()
	}
}

func init() {
	staticVersionsFlags = flag.NewFlagSet("static-versions", flag.ExitOnError)
	staticVersionsFlags.StringVar(&releaseParams.repo, "repo", ".", "path to the SDK git repository")
	staticVersionsFlags.BoolVar(&releaseParams.pretty, "pretty", false, "print indented JSON output")
	staticVersionsFlags.Var(&releaseParams.selector, "selector", "sets versioning strategy: release, development, or tags")
	staticVersionsFlags.Usage = func() {
		fmt.Printf("%s release static-versions\n", os.Args[0])
		staticVersionsFlags.PrintDefaults()
	}

	createReleaseFlags = flag.NewFlagSet("create", flag.ExitOnError)
	createReleaseFlags.StringVar(&releaseParams.repo, "repo", ".", "path to the SDK git repository")
	createReleaseFlags.StringVar(&releaseParams.releaseID, "id", "", "the ID of the release (e.g. 2020-07-17)")
	createReleaseFlags.BoolVar(&releaseParams.push, "push", false, "controls whether to push the release commit and tags to upstream repository")
	createReleaseFlags.BoolVar(&releaseParams.nonInteractive, "non-interactive", false, "bypass interactive prompts")
	createReleaseFlags.Usage = func() {
		fmt.Printf("%s release create [-repo=<repo>]\n", os.Args[0])
		createReleaseFlags.PrintDefaults()
	}
}

func releaseSubcmd(args []string) error {
	if len(args) == 0 {
		releaseUsage()
		return errors.New("invalid usage")
	}
	subCmd := args[0]

	switch subCmd {
	case "static-versions":
		err := staticVersionsFlags.Parse(args[1:])
		if err != nil {
			return err
		}

		return staticVersionsCmd(releaseParams.repo, releaseParams.selector, releaseParams.pretty)
	case "create":
		err := createReleaseFlags.Parse(args[1:])
		repo, err := changes.NewRepository(releaseParams.repo)
		if err != nil {
			return fmt.Errorf("couldn't load repository: %v", err)
		}

		err = repo.DoRelease(releaseParams.releaseID, releaseParams.push, !releaseParams.nonInteractive)
		if err != nil {
			return err
		}

		log.Printf("successfully created release %s\n", releaseParams.releaseID)

		return nil
	default:
		releaseUsage()
		return errors.New("invalid usage")
	}
}

func staticVersionsCmd(repoPath string, selector changes.VersionSelector, pretty bool) error {
	repo, err := changes.NewRepository(repoPath)
	if err != nil {
		return fmt.Errorf("couldn't load repository: %v", err)
	}

	enclosure, _, err := repo.DiscoverVersions(selector)
	if err != nil {
		return err
	}

	var out []byte
	if pretty {
		out, err = json.MarshalIndent(enclosure, "", "  ")
	} else {
		out, err = json.Marshal(enclosure)
	}
	if err != nil {
		return err
	}

	fmt.Print(string(out))
	return nil
}
