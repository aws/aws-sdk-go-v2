# Go SDK Module Versioning Tool
The `gosdk` versioning tool serves two purposes: allowing developers to annotate
the SDK modules with change metadata and allowing the SDK's automated release process
to consume this change metadata to manage module versions.

## Usage
The tool is split into two subcommands: change and release. The change subcommand is used
by developers to create metadata, and the release subcommand consumes this metadata.

### Change subcommand
The change subcommand allows a developer to manage change metadata within a repository. The command
must be run in a directory that is a descendant of a root directory that has a `.changes` directory.

Usage: `gosdk change [add,rm,edit,ls]`

The most common change operation is `gosdk change add`, which creates a new JSON change metadata
file under the `.changes/next-release` directory. The change subcommand is module aware and will
detect what module contains the working directory.

When creating change metadata, multiple modules can be entered into the editor template. This will
result in the creation of multiple change metadata files, one for each module.

Parts of the template may be passed to the command line tool via flags:
`gosdk change add --type feature --module my/module --description "description example"`. If all three
flags are provided, then the change will be created without opening the editor.

The `rm`, `edit`, and `ls` commands work similarly to `gosdk change add`, but allow you to delete, modify, and list
existing change metadata.

#### Wildcard changes
To create a change that applies to many modules, it is more convenient to use a wildcard change.
To create a wildcard change, the `--module` flag must be passed to `gosdk change add`. The provided module
uses a wildcard syntax similar to the `go` tooling. Addding a `/...` suffix to the end of a module
will match all modules in the repository with that prefix. For example:

`gosdk change add --module service/...` will match all modules with the prefix `service`.

Supplying only a wildcard module will match all modules with the given prefix regardless of whether they
have been changed or not. These modules will appear under an `affected_modules` field in the editor template
and can be manually edited to include only those modules that have changed. However, the `--compare-to` flag
provides a more convenient method of selecting only those modules that have actually been changed. In order to use
the `--compare-to` flag, you must first generate a snapshot of the repository before making changes in the form of
a version enclosure:

`gosdk release static-versions --selector release --repo . > ~/enc.json`

The `static-versions` command is described more below, but once we have generated `enc.json`, we can make changes
to the repository and then run the following command to create a wildcard change applying to only those module
that differ from the snapshot in `enc.json`:

`gosdk change add --module service/... --compare-to ~/enc.json`


### Release subcommand
The release subcommand currently serves two functions: generating a static version enclosure and 
creating a release. The `static-versions` command usage is:

`gosdk release static-versions --selector [tags, release, development] --repo [path]`

This command will print a snapshot of the module versions to `stdout`. The exact version associated
with each module is determined by the given `selector`:
* `tags` returns the latest tagged version for each module.
* `release` returns what the version of each module will be after accounting for change metadata present in `.changes/next-release`.
* `development` returns commit hash pseudo-versions for each module referencing the latest commit.
  
The version enclosure will also include module hashes for each module at the selected version. The `gosdk change add` 
command's `--compare-to` flag uses these hashes to resolve module differences.

The `gosdk release create --repo [path] --id [release-id]` performs a release of the SDK, for example:
`gosdk release create --repo . --id 2020-08-13` will perform the following steps:
* Determine the next version of each module based on change metadata in `.changes/next-release`
* Recursively resolve intra-SDK dependencies so that each module depends on the latest version of any other SDK module.
    * This step potentially causes patch version bumps to modules whose only change was to their dependencies, since we
    update the module's go.mod file.
* A release file `.changes/releases/2020-08-13.json` is created, containing all change metadata and version bumps contained
  in this release.
* Top level and per-module `CHANGELOG.md` files are updated with the new entry prepended to the top of the file.
  * These files will be created if they do not already exist.
* For each module whose version was bumped, if the module has a `version.go` file, the line containing the module's version
  constant is updated with its new version.
* The whole git repository is staged, committed, tagged with new versions, and pushed.

