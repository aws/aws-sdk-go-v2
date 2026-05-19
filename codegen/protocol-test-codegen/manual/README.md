# Protocol Test Codegen – Manual Files

This directory contains hand-written Go source files that are copied into
`internal/protocoltest/` after code generation. They supplement the
auto-generated protocol tests with cases that cannot be expressed through
Smithy protocol test traits.

## How protocol test generation works

The build is driven by `build.gradle.kts` in the parent directory. Running
`./gradlew :protocol-test-codegen:build` executes the following task chain:

1. **generate-smithy-build** – Loads all Smithy service models from the
   `smithy-aws-protocol-tests` / `smithy-protocol-tests` dependency JARs and
   the local `models/` directory. For each discovered service (minus those in
   `excludedServices`) it writes a projection entry into `smithy-build.json`.
   Projection and module names are derived automatically from the service shape
   ID unless an explicit override is defined in the `overrides` map.

2. **buildSdk** – Runs Smithy Build with the generated `smithy-build.json`,
   producing Go client code for every projection.

3. **cleanProtocolTests** – Deletes the entire `internal/protocoltest/`
   directory so stale files from previous runs don't persist.

4. **copyGoCodegen** – Copies the freshly generated Go code into
   `internal/protocoltest/`.

5. **copyManualFiles** – Copies everything under this `manual/` directory into
   `internal/protocoltest/`, overlaying the generated output.

## Adding a manual file

Mirror the target path relative to `internal/protocoltest/`. For example, to
add a file that should end up at
`internal/protocoltest/jsonrpc10/my_test.go`, place it at:

```
manual/jsonrpc10/my_test.go
```

No build script changes are needed — the `copyManualFiles` task copies the
entire `manual/` tree.