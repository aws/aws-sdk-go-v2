# Description

Detects new and changed Go modules in the repository, associates changelog annotations, and computes the next semver
version tag for each module. Produces a release manifest that is used with other utilities to orchestrate a release.

# Usage

```
calculaterelease [-o <manifestFile>]
```

# Determining Modules for Release

The `calculaterelease` traverses the repository to discover Go modules that are present. Using the discovered module
locations the tool determines whether a module is a new module, or an existing module with changes. To do so, `Git` is
first used to retrieve a list of tags for the repository. These tags are filtered and then sorted for each module by
using Go's [module versioning rules][modules-version-numbers]. For modules that have been previously tagged,
[git-diff-tree] is used to determine if changes were made to the module path by comparing the latest tag to the current
repository HEAD. If the module path, excluding child sub-module directories, contains changes to either `*.go` or
`go.mod` the module is considered as having source changes and will be considered for release. New modules that have not
been tagged previously are considered to eligible for release.

After determining the set of modules that are eligible for release, `calculaterelease` builds a reverse-dependency tree
graph to incrementally mark modules as requiring a version bump if one or more it's dependencies or
transitive-dependencies has been determined to be changed.

Finally, after determining the complete change set the next module version is chosen by using the change annotations
created using the [changelog] tool to refine and compute the next desired version. See
[here](#determining-the-next-module-version) for a more in-depth description about version selection. After computing
the next version for each module a final pass occurs to filter out modules that are configured to not be tagged for
released. After this is complete a summary manifest of the final set of modules to be released is printed out to
standard output, or the desired output file location.

**IMPORTANT**: `calculaterelease` does not consider the working directory or index when determining what should be
released, thus all relevant changes MUST be committed before executing the command.

# Determining the Next Module Version

Determining the next version of a module is determined by a combination of heuristics, and the information provided by
change annotations to refine the next version of a given module. The precedence of one or more annotations is as follows:
`release > feature > bugfix >= dependency >= documentation >= announcement`. the precedence ordering defines the type
of semantic version bump that will occur, allowing for multiple annotations to be defined safely.

The following table summarizes a complete set of examples of how module version selection works.

Module Path | Latest Tag | Next Tag | Annotations | Config | Descriptions
--- | --- | --- | --- | --- | ---
`github.com/aws/aws-sdk-go-v2/foo` | N/A | `foo/v1.0.0-preview` | N/A | N/A | New repository modules with no annotations default to a preview release tag
`github.com/aws/aws-sdk-go-v2/foo` | `foo/v1.0.0-preview` | `foo/v1.0.0-preview.1` | `feature, bugfix` | N/A | All changes to pre-release semver tags will increment an integer separated by a `.` on the pre-release tag.
`github.com/aws/aws-sdk-go-v2/foo` | `foo/v1.0.0-preview.1` | `foo/v1.0.0-rc` | `bugfix` | `{"pre_release": "rc"}` | The `pre_release` config for a module can be used to control the semver pre-release identifier. Annotations that are not `release` do not affect the version increment behavior.
`github.com/aws/aws-sdk-go-v2/foo` | `foo/v1.0.0-rc` | `foo/v1.0.0` | `release` | `{"pre_release": "rc"}` | A release annotation indicates the module's pre-release tag should be removed
`github.com/aws/aws-sdk-go-v2/foo` | `foo/v1.0.0` | `foo/v1.0.1` | N/A | N/A | Modules with changes but no annotations default to patch bump
`github.com/aws/aws-sdk-go-v2/foo` | `foo/v1.0.1` | `foo/v1.0.2` | `bugfix` | N/A | Modules with a bugfix annotation will increment the patch component.
`github.com/aws/aws-sdk-go-v2/foo` | `foo/v1.0.2` | `foo/v1.1.0` | `feature` | N/A | Feature bump will increment the minor version component.
`github.com/aws/aws-sdk-go-v2/foo` | `foo/v1.1.0` | `foo/v1.2.0-alpha` | N/A | `{"pre_release": "alpha"}` | The `pre_release` configuration can be used to mark the a modules next tagged release as a pre-release. Pre-release tags always imply a feature bump when calculating the preview version.
`github.com/aws/aws-sdk-go-v2/foo/v2` | N/A | `foo/v2.0.0-preview` | N/A | N/A | New repository modules with no annotations default to a preview release tag
`github.com/aws/aws-sdk-go-v2/bar` | N/A | N/A | `release` | N/A | New repository modules can be marked with `release` annotation to be immediately tagged with non-pre-release tag.
`github.com/aws/aws-sdk-go-v2/baz` | N/A | N/A | `feature` | `{"no_tag": true}` | Modules that are configured with`no_tag` will not be tagged regardless of whether there are Git changes or annotations.

# Understanding a Release Manifest

A [JSON Schema][json-schema] definition is available that provides a description of the release manifest produced by this tool.
You can view the definition [here](../../release/manifest_schema.json).

# Configuration

At the repository root one or more keys can be added to the `modules` dictionary in the `modman.toml`.

```toml
[modules."relative/mod/path"]
no_tag = false   # Set to true to indicate that the module should not be tagged for release. Regardless of changes, annoations, or having previously been tagged.
pre_release = "" # Set a semantic version pre-release identifer that will be used in the next release.
```

# Examples

## Preview changes to be released

By default `calculaterelease` will output the manifest to STDOUT, making it easy to preview changes that are currently
pending release.

```
$ calculaterelease
{
    "id": "2021-05-07",
    "modules": {
        "feature/s3/manager": {
            "module_path": "github.com/aws/aws-sdk-go-v2/feature/s3/manager",
            "from": "v1.1.4",
            "to": "v1.1.5",
            "changes": {
                "dependency_update": true
            }
        },
        "service/internal/s3shared": {
            "module_path": "github.com/aws/aws-sdk-go-v2/service/internal/s3shared",
            "from": "v1.2.3",
            "to": "v1.2.4",
            "changes": {
                "source_change": true
            }
        },
        "service/s3": {
            "module_path": "github.com/aws/aws-sdk-go-v2/service/s3",
            "from": "v1.6.0",
            "to": "v1.6.1",
            "changes": {
                "dependency_update": true
            }
        },
        "service/s3control": {
            "module_path": "github.com/aws/aws-sdk-go-v2/service/s3control",
            "from": "v1.5.1",
            "to": "v1.5.2",
            "changes": {
                "dependency_update": true
            }
        }
    },
    "tags": [
        "feature/s3/manager/v1.1.5",
        "service/internal/s3shared/v1.2.4",
        "service/s3/v1.6.1",
        "service/s3control/v1.5.2"
    ]
}
```

## Calculate Release and Write Computation Manifest to File

```
$ calculaterelease -o /output/file
```

[json-schema]: https://json-schema.org/

[changelog]: ../changelog/README.md

[modules-version-numbers]: https://golang.org/doc/modules/version-numbers

[git-diff-tree]: https://git-scm.com/docs/git-diff-tree
