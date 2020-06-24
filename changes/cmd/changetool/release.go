package main

import (
	"fmt"
	"github.com/aggagen/changes"
	"os"
)

func releaseSubcmd(args []string) {
	if len(args) < 2 {
		usage()
		os.Exit(1)
	}

	repo, err := changes.NewRepository(args[0])
	if err != nil {
		fmt.Printf("couldn't load repository: %v\n", err)
		os.Exit(1)
	}

	switch args[1] {
	case "update-pending":
		release, err := repo.Metadata.CreateRelease("pending", map[string]changes.VersionBump{}, true)
		if err != nil {
			fmt.Printf("couldn't update pending changelogs: %v\n", err)
			os.Exit(1)
		}

		err = repo.UpdateChangelog(release, true)
		if err != nil {
			fmt.Printf("failed to update CHANGELOG_PENDING: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("successfully updated CHANGELOG_PENDING")
	}
}
