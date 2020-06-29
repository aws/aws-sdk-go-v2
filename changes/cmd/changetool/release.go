package main

import (
	"fmt"
	"github.com/aggagen/changes"
	"log"
)

func releaseSubcmd(args []string) {
	if len(args) < 2 {
		usage()
	}

	repo, err := changes.NewRepository(args[0])
	if err != nil {
		log.Fatalf("couldn't load repository: %v", err)
	}

	switch args[1] {
	case "update-pending":
		updatePendingCmd(repo)
	case "demo-release":
		release, err := repo.Metadata.CreateRelease("2020-06-26", map[string]changes.VersionBump{
			"changes": {
				From: "v1.0.0",
				To:   "v1.0.1",
			},
			"test": {
				From: "v1.2.3",
				To:   "v1.3.0",
			},
		}, false)
		if err != nil {
			panic(err)
		}

		repo.UpdateChangelog(release, false)
	}
}

func updatePendingCmd(repo *changes.Repository) {
	err := repo.UpdatePendingChangelog()
	if err != nil {
		log.Fatalf("failed to update CHANGELOG_PENDING: %v", err)
	}

	fmt.Println("successfully updated CHANGELOG_PENDING")
}
