# AWS SDK for Go V2 Documentation

This source tree contains the AWS SDK for Go V2 developer and migration guides.

The documentation is generated using [Hugo](https://gohugo.io/), which is a static website generator. This project uses
the [Docsy](https://www.docsy.dev/) Hugo theme for presentation and formatting of the docs content type.

## Development

### Getting Started

1. [Install Hugo](https://gohugo.io/getting-started/installing)
1. Verify Hugo
   ```bash
    hugo version
    ```
   **Output**:
   ```
   Hugo Static Site Generator v0.78.0-FD62817B/extended darwin/amd64 BuildDate: 2020-11-03T13:20:38Z
   ```

   This project requires that the `extended` version of Hugo is installed. Follow the download instructions to ensure
   the correct version is present.
1. Fork the SDK [repository](https://github.com/aws/aws-sdk-go-v2)
1. Checkout the documentation branch replacing `USERNAME` with your GitHub username.
   ```bash
   git clone --single-branch -b documentation git@github.com:USERNAME/aws-sdk-go-v2.git aws-sdk-go-v2-docs
   cd aws-sdk-go-v2-docs
   ```
1. Initialize Project Submodules
   ```bash
   git submodule update --init --recursive themes/docsy
   ```
1. Install [NodeJS (LTS)](https://nodejs.org/en/)
1. Install Project Dependencies
   ```bash
   npm install
   ```

### Previewing Changes

Hugo comes with a built-in development server that allows you to iterate and preview your documentation changes in
realtime. This can be done by using the following Hugo command:

```bash
hugo server
```

### Submitting PRs
Pull requests should be submitted to the [SDK Repository](sdkrepo), to help speed up the process and reduce the time
to merge time please ensure that `Allow edits and access to secrets by maintainers` is checked before submitting your PR.
This will allow the project maintainers to make minor adjustments or improvements to the submitted PR, allow us to reduce the
roundtrip time for merging your request.

### Building Content for Release
To generate the documentation content in preparation for releasing the updates, the following command should be run.
```
hugo -d docs
```
This command will update the static version of the documentation into the `docs/` folder.

### References
* [Hugo Documentation](https://gohugo.io/documentation/)
* [Docsy Documentation](https://www.docsy.dev/docs/)

[sdkrepo]: https://github.com/aws/aws-sdk-go-v2
