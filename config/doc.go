// Package config provides utilities for loading configuration from multiple
// sources that can be used to configure the SDK's API clients, and utilities.
//
// The config package will load configuration from environment variables, AWS
// shared configuration file (~/.aws/config), and AWS shared credentials file
// (~/.aws/credentials).
//
// Use the LoadDefaultConfig to load configuration from all the SDK's supported
// sources, and resolve credentials using the SDK's default credential chain.
//
// LoadDefaultConfig allows for a variadic list of additional configuration sources that can
// provide one or more configuration values which can be used to programmatically control the resolution
// of a specific value, or allow for broader range of additional configuration sources not supported by the SDK.
// These configuration sources will take precedence over the default configuration sources loaded by the SDK.
// A number of simple helpers (prefixed by ``With``)  are provided that allow for such programmatic overriding
// of specific configuration values.
package config
