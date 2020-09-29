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
// * TODO Additional documentation needed.
package config
