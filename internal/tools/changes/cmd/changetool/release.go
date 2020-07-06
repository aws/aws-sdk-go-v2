package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/internal/tools/changes"
	"os"
)

var updatePendingFlags *flag.FlagSet

func releaseUsage() {
	var sets = []*flag.FlagSet{updatePendingFlags}

	for _, f := range sets {
		f.Usage()
	}
}

func init() {
	updatePendingFlags = flag.NewFlagSet("update-pending", flag.ExitOnError)
	updatePendingFlags.Usage = func() {
		fmt.Printf("%s release update-pending <repo>\n  <repo>: path to git repository\n", os.Args[0])
		updatePendingFlags.PrintDefaults()
	}
}

func releaseSubcmd(args []string) error {
	if len(args) < 2 {
		releaseUsage()
		return errors.New("invalid usage")
	}

	subCmd := args[0]
	repoPath := args[1]

	repo, err := changes.NewRepository(repoPath)
	if err != nil {
		return fmt.Errorf("couldn't load repository: %v", err)
	}

	switch subCmd {
	case "update-pending":
		err = updatePendingFlags.Parse(args[1:])
		if err != nil {
			return err
		}

		return updatePendingCmd(repo)
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
