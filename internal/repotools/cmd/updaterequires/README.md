# Usage

```
updaterequires [-force] [-release <manifestFile>]

Options:
-release <manifestFile> Uses next computed version tag information from a release manifest to update module dependencies.
-force                  Force can be used to allow the tool to downgrade a dependency to a lower version.
                        By default a dependency is only updated if the go.mod recorded version is semantically lower.
```

# Description

`updaterequires` uses the local repositories Git tags to update all inter-repository module dependencies to their latest
tags. Additionally. dependency configuration recorded in root repository directory's `modman.toml` will be used to
update external Go module dependencies to a specific version.

When executing the utility an additional `-release` argument can be provided by entering a path to a release manifest
(Computed using `calculaterelease`). When provided a release manifest, the tool will project the computed tags onto the
existing repository tags, allowing the tool to update Go Module dependencies with the latest tags being considered
available.

Lastly in the event that a dependency needs to be forced to a particular version that is lower than what is currently
recorded in the `go.mod`, the `-force` flag can be used. The force flag only applies to external dependencies, and
when enabled will update a dependency to the recorded version indicated in `modman.toml` regardless of the `go.mod`
recorded version being semantically higher or lower.
