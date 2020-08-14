/*
Package changes implements Go module version management for a multi-module git repository.

Repository

The Repository struct is a representation of a git repository containing multiple Go modules. Repository handles the
whole release process with respect to versioning, including tagging the git repository and updating CHANGELOG, version.go,
and go.mod files. The sequence of steps required to create a release of the SDK can be run from the Repository's
DoRelease function:

	repo, err := changes.NewRepository("path/to/repository")
	if err != nil {
		panic(err)
	}

	repo.DoRelease("2020-08-12")

After this code snippet runs, the versions of modules within the specified repository will be updated based on change metadata
contained in the repository.

Metadata

In order to make module versioning decisions during a release, the Repository's DoRelease function consumes metadata
about changes included in the release. The metadata is stored in the root of the repository in a .changes directory.
The Metadata struct provides a way to load, create, and modify the metadata in the .changes directory. The .changes directory
is structured as follows:
	.changes/
		next-release/	contains Change metadata for pending changes (to be included in the next release).
		releases/	contains Release metadata for previous releases.
		versions.json	is a VersionEnclosure representing the state of the repository at the last release.

There are three file formats used in the .changes directory corresponding to three structs: Change, Release, and
VersionEnclosure.

Changelogs

The Repository's DoRelease function pdates both a top level CHANGELOG.md and per-module CHANGELOG.md files. The CHANGELOG.md
in the root of the repository provides a consolidated view of all changelog entries for every module. This means that the
root module does not have its own CHANGELOG.

The Release struct handles the rendering of CHANGELOG files, generating an entry that is prepended to CHANGELOG files
during a release.

Dependencies within the SDK

As part of the Repository's release process, any module within the repository that depends upon another module in the
repository will have its dependency updated to the latest version. So, after a release all modules in the SDK depend on
the latest version of any other SDK module.
*/
package changes
