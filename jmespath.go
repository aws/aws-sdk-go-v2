package sdk

import "github.com/jmespath/go-jmespath"

// Temporarily pin go-jmespath as a direct dependency.
// FUTURE: remove this once all waiters are code-generated.
var _ = jmespath.Search
