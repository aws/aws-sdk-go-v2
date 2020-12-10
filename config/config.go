package config

import (
	"context"
	"io"

	"github.com/awslabs/smithy-go/logging"
	"github.com/awslabs/smithy-go/middleware"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go-v2/credentials/endpointcreds"
	"github.com/aws/aws-sdk-go-v2/credentials/processcreds"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
)

// defaultLoaders are a slice of functions that will read external configuration
// sources for configuration values. These values are read by the AWSConfigResolvers
// using interfaces to extract specific information from the external configuration.
var defaultLoaders = []loader{
	loadEnvConfig,
	loadSharedConfigIgnoreNotExist,
}

// defaultAWSConfigResolvers are a slice of functions that will resolve external
// configuration values into AWS configuration values.
//
// This will setup the AWS configuration's Region,
var defaultAWSConfigResolvers = []awsConfigResolver{
	// Resolves the default configuration the SDK's aws.Config will be
	// initialized with.
	resolveDefaultAWSConfig,

	// Sets the logger to be used. Could be user provided logger, and client
	// logging mode.
	resolveLogger,
	resolveClientLogMode,

	// Sets the HTTP client and configuration to use for making requests using
	// the HTTP transport.
	resolveHTTPClient,
	resolveCustomCABundle,

	// Sets the endpoint resolving behavior the API Clients will use for making
	// requests to. Clients default to their own clients this allows overrides
	// to be specified.
	resolveEndpointResolver,

	// Sets the retry behavior API clients will use within their retry attempt
	// middleware. Defaults to unset, allowing API clients to define their own
	// retry behavior.
	resolveRetryer,

	// Sets the region the API Clients should use for making requests to.
	resolveRegion,
	// TODO: Add back EC2 Region Resolver Support
	resolveDefaultRegion,

	// Sets the additional set of middleware stack mutators that will custom
	// API client request pipeline middleware.
	resolveAPIOptions,

	// Sets the resolved credentials the API clients will use for
	// authentication. Provides the SDK's default credential chain.
	//
	// Should probably be the last step in the resolve chain to ensure that all
	// other configurations are resolved first in case downstream credentials
	// implementations depend on or can be configured with earlier resolved
	// configuration options.
	resolveCredentials,
}

// A Config represents a generic configuration value or set of values. This type
// will be used by the AWSConfigResolvers to extract
//
// General the Config type will use type assertion against the Provider interfaces
// to extract specific data from the Config.
type Config interface{}

// A loader is used to load external configuration data and returns it as
// a generic Config type.
//
// The loader should return an error if it fails to load the external configuration
// or the configuration data is malformed, or required components missing.
type loader func(context.Context, configs) (Config, error)

// An awsConfigResolver will extract configuration data from the configs slice
// using the provider interfaces to extract specific functionality. The extracted
// configuration values will be written to the AWS Config value.
//
// The resolver should return an error if it it fails to extract the data, the
// data is malformed, or incomplete.
type awsConfigResolver func(ctx context.Context, cfg *aws.Config, configs configs) error

// configs is a slice of Config values. These values will be used by the
// AWSConfigResolvers to extract external configuration values to populate the
// AWS Config type.
//
// Use AppendFromLoaders to add additional external Config values that are
// loaded from external sources.
//
// Use ResolveAWSConfig after external Config values have been added or loaded
// to extract the loaded configuration values into the AWS Config.
type configs []Config

// AppendFromLoaders iterates over the slice of loaders passed in calling each
// loader function in order. The external config value returned by the loader
// will be added to the returned configs slice.
//
// If a loader returns an error this method will stop iterating and return
// that error.
func (cs configs) AppendFromLoaders(ctx context.Context, loaders []loader) (configs, error) {
	for _, fn := range loaders {
		cfg, err := fn(ctx, cs)
		if err != nil {
			return nil, err
		}

		cs = append(cs, cfg)
	}

	return cs, nil
}

// ResolveAWSConfig returns a AWS configuration populated with values by calling
// the resolvers slice passed in. Each resolver is called in order. Any resolver
// may overwrite the AWS Configuration value of a previous resolver.
//
// If an resolver returns an error this method will return that error, and stop
// iterating over the resolvers.
func (cs configs) ResolveAWSConfig(ctx context.Context, resolvers []awsConfigResolver) (aws.Config, error) {
	var cfg aws.Config

	for _, fn := range resolvers {
		if err := fn(ctx, &cfg, cs); err != nil {
			// TODO provide better error?
			return aws.Config{}, err
		}
	}

	var sources []interface{}
	for _, s := range cs {
		sources = append(sources, s)
	}
	cfg.ConfigSources = sources

	return cfg, nil
}

// ResolveConfig calls the provide function passing slice of configuration sources.
// This implements the aws.ConfigResolver interface.
func (cs configs) ResolveConfig(f func(configs []interface{}) error) error {
	var cfgs []interface{}
	for i := range cs {
		cfgs = append(cfgs, cs[i])
	}
	return f(cfgs)
}

// LoadOptionsFunc is a type alias for LoadOptions functional option
type LoadOptionsFunc func(*LoadOptions) error

// LoadOptions are discrete set of options that are valid for loading the configuration
type LoadOptions struct {

	// Region is the region to send requests to.
	Region string

	// Credentials object to use when signing requests.
	Credentials aws.CredentialsProvider

	// HTTPClient the SDK's API clients will use to invoke HTTP requests.
	HTTPClient HTTPClient

	// EndpointResolver that can be used to provide or override an endpoint for the given
	// service and region Please see the `aws.EndpointResolver` documentation on usage.
	EndpointResolver aws.EndpointResolver

	// Retryer guides how HTTP requests should be retried in case of
	// recoverable failures.
	Retryer aws.Retryer

	// ConfigSources are the sources that were used to construct the Config.
	// Allows for additional configuration to be loaded by clients.
	ConfigSources []interface{}

	// APIOptions provides the set of middleware mutations modify how the API
	// client requests will be handled. This is useful for adding additional
	// tracing data to a request, or changing behavior of the SDK's client.
	APIOptions []func(*middleware.Stack) error

	// Logger writer interface to write logging messages to.
	Logger logging.Logger

	// ClientLogMode is used to configure the events that will be sent to the configured logger.
	// This can be used to configure the logging of signing, retries, request, and responses
	// of the SDK clients.
	//
	// See the ClientLogMode type documentation for the complete set of logging modes and available
	// configuration.
	ClientLogMode *aws.ClientLogMode

	// SharedConfigProfile is the profile to be used when loading the SharedConfig
	SharedConfigProfile string

	// SharedConfigFiles is the slice of custom shared config files to use when loading the SharedConfig
	SharedConfigFiles []string

	// CustomCABundle is CA bundle PEM bytes reader
	CustomCABundle io.Reader

	// DefaultRegion is the fall back region, used if a region was not resolved from other sources
	DefaultRegion string

	// UseEC2IMDSRegion indicates if SDK should retrieve the region
	// from the EC2 Metadata service
	UseEC2IMDSRegion *UseEC2IMDSRegion

	// ProcessCredentialOptions is a function for setting
	// the processcreds.Options
	ProcessCredentialOptions func(*processcreds.Options)

	// EC2RoleCredentialOptions is a function for setting
	// the ec2rolecreds.Options
	EC2RoleCredentialOptions func(*ec2rolecreds.Options)

	// EndpointCredentialOptions is a function for setting
	// the endpointcreds.Options
	EndpointCredentialOptions func(*endpointcreds.Options)

	// WebIdentityRoleCredentialOptions is a function for setting
	// the stscreds.WebIdentityRoleOptions
	WebIdentityRoleCredentialOptions func(*stscreds.WebIdentityRoleOptions)

	// AssumeRoleCredentialOptions is a function for setting the
	// stscreds.AssumeRoleOptions
	AssumeRoleCredentialOptions func(*stscreds.AssumeRoleOptions)

	//LogConfigurationWarnings when set to true, enables logging
	// configuration warnings
	LogConfigurationWarnings *bool
}

// LoadDefaultConfig reads the SDK's default external configurations, and
// populates an AWS Config with the values from the external configurations.
//
// An optional variadic set of additional Config values can be provided as input
// that will be prepended to the configs slice. Use this to add custom configuration.
// The custom configurations must satisfy the respective providers for their data
// or the custom data will be ignored by the resolvers and config loaders.
//
//    cfg, err := config.LoadDefaultConfig( context.TODO(),
//       WithSharedConfigProfile("test-profile"),
//    )
//    if err != nil {
//       panic(fmt.Sprintf("failed loading config, %v", err))
//    }
//
//
// The default configuration sources are:
// * Environment Variables
// * Shared Configuration and Shared Credentials files.
func LoadDefaultConfig(ctx context.Context, optFns ...func(*LoadOptions) error) (cfg aws.Config, err error) {
	var options LoadOptions
	for _, optFn := range optFns {
		optFn(&options)
	}

	// assign Load Options to configs
	var cfgCpy = configs{options}

	cfgCpy, err = cfgCpy.AppendFromLoaders(ctx, defaultLoaders)
	if err != nil {
		return aws.Config{}, err
	}

	cfg, err = cfgCpy.ResolveAWSConfig(ctx, defaultAWSConfigResolvers)
	if err != nil {
		return aws.Config{}, err
	}

	return cfg, nil
}

// getRegion returns Region from config's LoadOptions
func (o LoadOptions) getRegion(ctx context.Context) (string, bool, error) {
	if len(o.Region) == 0 {
		return "", false, nil
	}

	return o.Region, true, nil
}

// WithRegion is a helper function to construct functional options
// that sets Region on config's LoadOptions. Setting the region to
// an empty string, will result in the region value being ignored.
func WithRegion(v string) LoadOptionsFunc {
	return func(o *LoadOptions) error {
		o.Region = v
		return nil
	}
}

// getDefaultRegion returns DefaultRegion from config's LoadOptions
func (o LoadOptions) getDefaultRegion(ctx context.Context) (string, bool, error) {
	if len(o.DefaultRegion) == 0 {
		return "", false, nil
	}

	return o.DefaultRegion, true, nil
}

// WithDefaultRegion is a helper function to construct functional options
// that sets a DefaultRegion on config's LoadOptions. Setting the default
// region to an empty string, will result in the default region value
// being ignored.
func WithDefaultRegion(v string) LoadOptionsFunc {
	return func(o *LoadOptions) error {
		o.DefaultRegion = v
		return nil
	}
}

// getSharedConfigProfile returns SharedConfigProfile from config's LoadOptions
func (o LoadOptions) getSharedConfigProfile(ctx context.Context) (string, bool, error) {
	if len(o.SharedConfigProfile) == 0 {
		return "", false, nil
	}

	return o.SharedConfigProfile, true, nil
}

// WithSharedConfigProfile is a helper function to construct functional options
// that sets SharedConfigProfile on config's LoadOptions. Setting the shared
// config profile to an empty string, will result in the shared config profile
// value being ignored.
func WithSharedConfigProfile(v string) LoadOptionsFunc {
	return func(o *LoadOptions) error {
		o.SharedConfigProfile = v
		return nil
	}
}

// getSharedConfigFiles returns SharedConfigFiles set on config's LoadOptions
func (o LoadOptions) getSharedConfigFiles(ctx context.Context) ([]string, bool, error) {
	if o.SharedConfigFiles == nil {
		return nil, false, nil
	}

	return o.SharedConfigFiles, true, nil
}

// WithSharedConfigFiles is a helper function to construct functional options
// that sets slice of SharedConfigFiles on config's LoadOptions. Setting the shared config
// files to an nil string slice, will result in the shared config files value being ignored.
func WithSharedConfigFiles(v []string) LoadOptionsFunc {
	return func(o *LoadOptions) error {
		o.SharedConfigFiles = v
		return nil
	}
}

// getCustomCABundle returns CustomCABundle from LoadOptions
func (o LoadOptions) getCustomCABundle(ctx context.Context) (io.Reader, bool, error) {
	if o.CustomCABundle == nil {
		return nil, false, nil
	}

	return o.CustomCABundle, true, nil
}

// WithCustomCABundle is a helper function to construct functional options
// that sets CustomCABundle on config's LoadOptions. Setting the custom CA Bundle
// to nil will result in custom CA Bundle value being ignored.
func WithCustomCABundle(v io.Reader) LoadOptionsFunc {
	return func(o *LoadOptions) error {
		o.CustomCABundle = v
		return nil
	}
}

// WithEC2IMDSRegion is a helper function to construct functional options
// that enables resolving EC2IMDS region. The function takes
// in a UseEC2IMDSRegion functional option, and can be used to set the EC2IMDS client
// which will be used to resolve EC2IMDSRegion. If no functional option is provided,
// an EC2IMDS client is built and used by the resolver.
func WithEC2IMDSRegion(fnOpts ...func(o *UseEC2IMDSRegion)) LoadOptionsFunc {
	return func(o *LoadOptions) error {
		o.UseEC2IMDSRegion = &UseEC2IMDSRegion{}

		for _, fn := range fnOpts {
			fn(o.UseEC2IMDSRegion)
		}
		return nil
	}
}

// getCredentialsProvider returns the credentials value
func (o LoadOptions) getCredentialsProvider(ctx context.Context) (aws.CredentialsProvider, bool, error) {
	if o.Credentials == nil {
		return nil, false, nil
	}

	return o.Credentials, true, nil
}

// WithCredentialsProvider is a helper function to construct functional options
// that sets Credential provider value on config's LoadOptions. If credentials provider
// is set to nil, the credentials provider value will be ignored.
func WithCredentialsProvider(v aws.CredentialsProvider) LoadOptionsFunc {
	return func(o *LoadOptions) error {
		o.Credentials = v
		return nil
	}
}

// getProcessCredentialOptions returns the wrapped function to set processcreds.Options
func (o LoadOptions) getProcessCredentialOptions(ctx context.Context) (func(*processcreds.Options), bool, error) {
	if o.ProcessCredentialOptions == nil {
		return nil, false, nil
	}

	return o.ProcessCredentialOptions, true, nil
}

// WithProcessCredentialOptions is a helper function to construct functional options
// that sets a function to use processcreds.Options on config's LoadOptions. If process
// credential options is set to nil, the process credential value will be ignored.
func WithProcessCredentialOptions(v func(*processcreds.Options)) LoadOptionsFunc {
	return func(o *LoadOptions) error {
		o.ProcessCredentialOptions = v
		return nil
	}
}

// getEC2RoleCredentialOptions returns the wrapped function to set the ec2rolecreds.Options
func (o LoadOptions) getEC2RoleCredentialOptions(ctx context.Context) (func(*ec2rolecreds.Options), bool, error) {
	if o.EC2RoleCredentialOptions == nil {
		return nil, false, nil
	}

	return o.EC2RoleCredentialOptions, true, nil
}

// WithEC2RoleCredentialOptions is a helper function to construct functional options
// that sets a function to use ec2rolecreds.Options on config's LoadOptions. If
// EC2 role credential options is set to nil, the EC2 role credential options value
// will be ignored.
func WithEC2RoleCredentialOptions(v func(*ec2rolecreds.Options)) LoadOptionsFunc {
	return func(o *LoadOptions) error {
		o.EC2RoleCredentialOptions = v
		return nil
	}
}

// getEndpointCredentialOptions returns the wrapped function to set endpointcreds.Options
func (o LoadOptions) getEndpointCredentialOptions(context.Context) (func(*endpointcreds.Options), bool, error) {
	if o.EndpointCredentialOptions == nil {
		return nil, false, nil
	}

	return o.EndpointCredentialOptions, true, nil
}

// WithEndpointCredentialOptions is a helper function to construct functional options
// that sets a function to use endpointcreds.Options on config's LoadOptions. If
// endpoint credential options is set to nil, the endpoint credential options value will
// be ignored.
func WithEndpointCredentialOptions(v func(*endpointcreds.Options)) LoadOptionsFunc {
	return func(o *LoadOptions) error {
		o.EndpointCredentialOptions = v
		return nil
	}
}

// getWebIdentityRoleCredentialOptions returns the wrapped function
func (o LoadOptions) getWebIdentityRoleCredentialOptions(context.Context) (func(*stscreds.WebIdentityRoleOptions), bool, error) {
	if o.WebIdentityRoleCredentialOptions == nil {
		return nil, false, nil
	}

	return o.WebIdentityRoleCredentialOptions, true, nil
}

// WithWebIdentityRoleCredentialOptions is a helper function to construct functional options
// that sets a function to use stscreds.WebIdentityRoleOptions on config's LoadOptions. If
// web identity role credentials options is set to nil, the web identity role credentials value
// will be ignored.
func WithWebIdentityRoleCredentialOptions(v func(*stscreds.WebIdentityRoleOptions)) LoadOptionsFunc {
	return func(o *LoadOptions) error {
		o.WebIdentityRoleCredentialOptions = v
		return nil
	}
}

// getAssumeRoleCredentialOptions returns AssumeRoleCredentialOptions from LoadOptions
func (o LoadOptions) getAssumeRoleCredentialOptions(context.Context) (func(options *stscreds.AssumeRoleOptions), bool, error) {
	if o.AssumeRoleCredentialOptions == nil {
		return nil, false, nil
	}

	return o.AssumeRoleCredentialOptions, true, nil
}

// WithAssumeRoleCredentialOptions  is a helper function to construct functional options
// that sets a function to use stscreds.AssumeRoleOptions on config's LoadOptions. If
// assume role credentials options is set to nil, the assume role credentials value
// will be ignored.
func WithAssumeRoleCredentialOptions(v func(*stscreds.AssumeRoleOptions)) LoadOptionsFunc {
	return func(o *LoadOptions) error {
		o.AssumeRoleCredentialOptions = v
		return nil
	}
}

func (o LoadOptions) getHTTPClient(ctx context.Context) (HTTPClient, bool, error) {
	if o.HTTPClient == nil {
		return nil, false, nil
	}

	return o.HTTPClient, true, nil
}

// WithHTTPClient is a helper function to construct functional options
// that sets HTTPClient on LoadOptions. If HTTPClient is set to nil,
// the HTTPClient value will be ignored.
func WithHTTPClient(v HTTPClient) LoadOptionsFunc {
	return func(o *LoadOptions) error {
		o.HTTPClient = v
		return nil
	}
}

func (o LoadOptions) getAPIOptions(ctx context.Context) ([]func(*middleware.Stack) error, bool, error) {
	if o.APIOptions == nil {
		return nil, false, nil
	}

	return o.APIOptions, true, nil
}

// WithAPIOptions is a helper function to construct functional options
// that sets APIOptions on LoadOptions. If APIOptions is set to nil, the
// APIOptions value is ignored.
func WithAPIOptions(v []func(*middleware.Stack) error) LoadOptionsFunc {
	return func(o *LoadOptions) error {
		if v == nil {
			return nil
		}

		o.APIOptions = append(o.APIOptions, v...)
		return nil
	}
}

func (o LoadOptions) getRetryer(ctx context.Context) (aws.Retryer, bool, error) {
	if o.Retryer == nil {
		return nil, false, nil
	}

	return o.Retryer, true, nil
}

// WithRetryer is a helper function to construct functional options
// that sets Retryer on LoadOptions. If Retryer is set to nil, the
// Retryer value is ignored.
func WithRetryer(v aws.Retryer) LoadOptionsFunc {
	return func(o *LoadOptions) error {
		o.Retryer = v
		return nil
	}
}

func (o LoadOptions) getEndpointResolver(ctx context.Context) (aws.EndpointResolver, bool, error) {
	if o.EndpointResolver == nil {
		return nil, false, nil
	}

	return o.EndpointResolver, true, nil
}

// WithEndpointResolver is a helper function to construct functional options
// that sets endpoint resolver on LoadOptions. The EndpointResolver is set to nil,
// the EndpointResolver value is ignored.
func WithEndpointResolver(v aws.EndpointResolver) LoadOptionsFunc {
	return func(o *LoadOptions) error {
		o.EndpointResolver = v
		return nil
	}
}

func (o LoadOptions) getLogger(ctx context.Context) (logging.Logger, bool, error) {
	if o.Logger == nil {
		return nil, false, nil
	}

	return o.Logger, true, nil
}

// WithLogger is a helper function to construct functional options
// that sets Logger on LoadOptions. If Logger is set to nil, the
// Logger value will be ignored.
func WithLogger(v logging.Logger) LoadOptionsFunc {
	return func(o *LoadOptions) error {
		o.Logger = v
		return nil
	}
}

func (o LoadOptions) getClientLogMode(ctx context.Context) (aws.ClientLogMode, bool, error) {
	if o.ClientLogMode == nil {
		return 0, false, nil
	}

	return *o.ClientLogMode, true, nil
}

// WithClientLogMode is a helper function to construct functional options
// that sets client log mode on LoadOptions. If client log mode is set to nil,
// the client log mode value will be ignored.
func WithClientLogMode(v aws.ClientLogMode) LoadOptionsFunc {
	return func(o *LoadOptions) error {
		o.ClientLogMode = &v
		return nil
	}
}

func (o LoadOptions) getLogConfigurationWarnings(ctx context.Context) (v bool, found bool, err error) {
	if o.LogConfigurationWarnings == nil {
		return false, false, nil
	}
	return *o.LogConfigurationWarnings, true, nil
}

// WithLogConfigurationWarnings is a helper function to construct functional options
// that can be used to set Log configuration warnings. If log configuration warnings
// is set to nil, the log configuration warnings value is ignored.
func WithLogConfigurationWarnings(v bool) LoadOptionsFunc {
	return func(o *LoadOptions) error {
		o.LogConfigurationWarnings = &v
		return nil
	}
}
