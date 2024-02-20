# v1.5.1 (2024-02-20)

* **Bug Fix**: When sourcing values for a service's `EndpointParameters`, the lack of a configured region (i.e. `options.Region == ""`) will now translate to a `nil` value for `EndpointParameters.Region` instead of a pointer to the empty string `""`. This will result in a much more explicit error when calling an operation instead of an obscure hostname lookup failure.

# v1.5.0 (2024-02-16)

* **Feature**: Add new ClientOptions field to waiter config which allows you to extend the config for operation calls made by waiters.

# v1.4.0 (2024-02-13)

* **Feature**: Bump minimum Go version to 1.20 per our language support policy.
* **Dependency Update**: Updated to the latest SDK module versions

# v1.3.0 (2024-02-12)

* **Feature**: Adding a new option "parameters" for data plane api ExecuteQuery to support running parameterized query via SDK.

# v1.2.0 (2024-02-01)

* **Feature**: Adding new APIs in SDK for Amazon Neptune Analytics. These APIs include operations to execute, cancel, list queries and get the graph summary.

# v1.1.1 (2024-01-04)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.1.0 (2023-12-21)

* **Feature**: Adds Waiters for successful creation and deletion of Graph, Graph Snapshot, Import Task and Private Endpoints for Neptune Analytics

# v1.0.0 (2023-12-14)

* **Release**: New AWS service client module
* **Feature**: This is the initial SDK release for Amazon Neptune Analytics

