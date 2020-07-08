package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/internal/tools/changes"
)

func releaseSubcmd(args []string) error {
	if len(args) < 2 {
		usage()
	}

	repo, err := changes.NewRepository(args[0])
	if err != nil {
		return fmt.Errorf("couldn't load repository: %v", err)
	}

	switch args[1] {
	case "update-pending":
		return updatePendingCmd(repo)
	case "versions-init":
		return repo.InitializeVersions()
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

		return repo.UpdateChangelog(release, false)
	case "test":
		enclosure, _ := repo.DiscoverVersions(changes.ReleaseVersionSelector)
		repo.Metadata.SaveEnclosure(enclosure)
	case "scratch":

	default:
		usage()
	}

	return nil
}

func updatePendingCmd(repo *changes.Repository) error {
	err := repo.UpdatePendingChangelog()
	if err != nil {
		return fmt.Errorf("failed to update CHANGELOG_PENDING: %v", err)
	}

	fmt.Println("successfully updated CHANGELOG_PENDING")
	return nil
}
