# v1.0.3 (2025-12-02)

* **Dependency Update**: Updated to the latest SDK module versions
* **Dependency Update**: Upgrade to smithy-go v1.24.0. Notably this version of the library reduces the allocation footprint of the middleware system. We observe a ~10% reduction in allocations per SDK call with this change.

# v1.0.2 (2025-11-25)

* **Bug Fix**: Add error check for endpoint param binding during auth scheme resolution to fix panic reported in #3234

# v1.0.1 (2025-11-19.2)

* **Dependency Update**: Updated to the latest SDK module versions

# v1.0.0 (2025-11-17)

* **Release**: New AWS service client module
* **Feature**: Amazon MWAA now offers serverless deployment, eliminating operational overhead while optimizing costs. The service supports YAML and Python-based workflows, with 80+ AWS Operators. It provides isolated execution, IAM permissions, and automatic scaling with pay-per-use pricing.

