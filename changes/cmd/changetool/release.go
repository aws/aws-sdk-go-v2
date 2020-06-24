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
	}
}

func updatePendingCmd(repo *changes.Repository) {
	err := repo.UpdatePendingChangelog()
	if err != nil {
		log.Fatalf("failed to update CHANGELOG_PENDING: %v", err)
	}

	fmt.Println("successfully updated CHANGELOG_PENDING")
}
