# Repository Management Tools

This Go module is a collection of tools that have been written for managing the AWS SDK for Go V2 source repository.
At the time these were written, the Go ecosystem lacked tooling for managing repositories containing multiple Go modules
with at the size and scale of the AWS SDK. With over 274 Go modules in the repository the tooling found here has been
made to manage the lifecycle of managing dependencies (iter-repository and external), managing the releases of new
changes, tag management following Go semver requirements, and production of changelogs.

## Utilities
The following is a breakdown of some the key utilities found in this module that are used to manage the SDK and handle
the complex release process.

Commands | Description | README
--- | --- | ---
`changelog` | Create and manage changelog annotations. Annotations are used to document module changes and refining of the next semver version. | [Link][changelog]
`updaterequires` | Manages `go.mod` require entries, allows for easily updating inter-repository module dependencies to their latest tag, and the ability to quickly manage external dependency requirements. | N/A
`updatemodulemeta` | Generates a `go_module_metadata.go` file in each module containing useful runtime metadata like the modules tagged version. | N/A
`generatechangelog` | Uses a release description and associated changelog annotations to produce `CHANGELOG.md` entries for the release in each repository module. In addition, a summarized release statement will be created at the root of the repository. | N/A
`gomodgen` | Copies [smithy-go] codegen build artifacts into the SDK repository and generates a `go.mod` file using the build artifacts `generated.json` description. | N/A
`annotatestablegen` | Generates a release changelog annotation type for **new** [smithy-go] generated modules that are not marked as unstable. | N/A
`calculaterelease` | Detects new and changed Go modules in the repository, associates changelog annotations, and computes the next semver version tag for each module. Produces a release manifest that is used with other utilities to orchestrate a release. | [Link][calculaterelease]
`tagrelease` | Commits pending changes to the working directory, reads the release manifest, and creates the computed tags | N/A
`makerelative` | Used to generate `go.mod` `replace` statements for inter-repository module dependencies. This ensures that when developing on a given Go module it's iter-repository dependencies refer to the cloned repository. | N/A
`eachmodule` | Utility for quickly scripting execution of commands in each module of a repository. | N/A

# Configuration

A number of the repository tools, specifically those involved with the dependency management and release have specific
behavior that is driven by the `modman.toml` file found at the root of the git repository. This configuration file is
a [TOML] configuration file.

## Dependencies
The `dependencies` is a dictionary of key-value pairs that describe **external** dependencies that one or more modules
within the repository may depend on. (External dependencies is defined as the set of Go modules that are not found 
within the project git repository.) This section is used to quickly set the version of a dependency modules in the
repository should use. The `updaterequires` tool can be used to update all Go modules require statements for each module
in the repository and update them to the indicated version if they depend on the given external module.

### Example
```toml
[dependencies]
"github.com/aws/smithy-go" = "v1.4.0"
```

This example indicates that repository modules that depend on `github.com/aws/smithy-go` should depend on `v1.4.0`
version of the library. After updating the value in the configuration file, `updaterequires` can be used to update
modules with this information.

## Modules

`modules` is a dictionary where the keys are module directories relative to the repository root. Each key maps to a
dictionary of key-value pairs that can configure several properties of a module that affect how or
if a module is handled when performing a release.

### Example
To configure the module `feature/service/shinything` to not be tagged by the release process:

```toml
[modules."feature/service/shinything"]
no_tag = true
```

For more information on how to configure how modules are released see the [calculaterelease README][calculaterelease].

#### Mark a module to be 

**NOTE**: If you wish to create a configuration item for a module located at the root of the repository use
`.` as the key name.

[calculaterelease]: cmd/calculaterelease/README.md
[changelog]: cmd/changelog/README.md
[smithy-go]: https://github.com/aws/smithy-go
[TOML]: https://toml.io
