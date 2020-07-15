package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/internal/tools/changes"
	"os"
)

var releaseParams = struct {
	repo     string
	selector changes.VersionSelector
}{}

var updatePendingFlags *flag.FlagSet
var staticVersionsFlags *flag.FlagSet

func releaseUsage() {
	var sets = []*flag.FlagSet{updatePendingFlags, staticVersionsFlags}

	for _, f := range sets {
		f.Usage()
	}
}

func init() {
	updatePendingFlags = flag.NewFlagSet("update-pending", flag.ExitOnError)
	updatePendingFlags.StringVar(&releaseParams.repo, "repo", "", "path to the SDK git repository")
	updatePendingFlags.Usage = func() {
		fmt.Printf("%s release update-pending <repo>\n", os.Args[0])
		updatePendingFlags.PrintDefaults()
	}

	staticVersionsFlags = flag.NewFlagSet("static-versions", flag.ExitOnError)
	staticVersionsFlags.StringVar(&releaseParams.repo, "repo", "", "path to the SDK git repository")
	staticVersionsFlags.Var(&releaseParams.selector, "selector", "sets versioning strategy: release, development, or tags")
	staticVersionsFlags.Usage = func() {
		fmt.Printf("%s release static-versions\n", os.Args[0])
		staticVersionsFlags.PrintDefaults()
	}
}

func releaseSubcmd(args []string) error {
	if len(args) == 0 {
		releaseUsage()
		return errors.New("invalid usage")
	}
	subCmd := args[0]

	switch subCmd {
	case "update-pending":
		err := updatePendingFlags.Parse(args[1:])
		if err != nil {
			return err
		}

		repo, err := changes.NewRepository(releaseParams.repo)
		if err != nil {
			return fmt.Errorf("couldn't load repository: %v", err)
		}

		return updatePendingCmd(repo)
	case "static-versions":
		err := staticVersionsFlags.Parse(args[1:])
		if err != nil {
			return err
		}

		return staticVersionsCmd(releaseParams.repo, releaseParams.selector)
	case "create":
		repo, err := changes.NewRepository(args[1])
		if err != nil {
			return fmt.Errorf("couldn't load repository: %v", err)
		}

		fmt.Println(repo.DoRelease())

		return nil
	case "test":

		return nil
	default:
		releaseUsage()
		return errors.New("invalid usage")
	}
}

func updatePendingCmd(repo *changes.Repository) error {
	err := repo.UpdatePendingChangelog()
	if err != nil {
		return fmt.Errorf("failed to update CHANGELOG_PENDING: %v", err)
	}

	fmt.Println("successfully updated CHANGELOG_PENDING")
	return nil
}

func staticVersionsCmd(repoPath string, selector changes.VersionSelector) error {
	repo, err := changes.NewRepository(repoPath)
	if err != nil {
		return fmt.Errorf("couldn't load repository: %v", err)
	}

	enclosure, _, err := repo.DiscoverVersions(selector)
	if err != nil {
		return err
	}

	out, err := json.Marshal(enclosure)
	if err != nil {
		return err
	}

	fmt.Print(string(out))
	return nil
}
