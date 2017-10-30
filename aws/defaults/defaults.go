// Package defaults is a collection of helpers to retrieve the SDK's default
// configuration and handlers.
//
// Generally this package shouldn't be used directly, but session.Session
// instead. This package is useful when you need to reset the defaults
// of a session or service client to the SDK defaults before setting
// additional parameters.
//
// TODO rename to "default"
package defaults

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
)

// Logger returns a Logger which will write log messages to stdout, and
// use same formatting runes as the stdlib log.Logger
func Logger() aws.Logger {
	return &defaultLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// A defaultLogger provides a minimalistic logger satisfying the Logger interface.
type defaultLogger struct {
	logger *log.Logger
}

// Log logs the parameters to the stdlib logger. See log.Println.
func (l defaultLogger) Log(args ...interface{}) {
	l.logger.Println(args...)
}

// Config returns the default configuration without credentials.
// To retrieve a config with credentials also included use
// `defaults.Get().Config` instead.
//
// Generally you shouldn't need to use this method directly, but
// is available if you need to reset the configuration of an
// existing service client or session.
func Config() aws.Config {
	return aws.Config{
		EndpointResolver: endpoints.DefaultResolver(),
		Credentials:      aws.AnonymousCredentials,
		HTTPClient:       HTTPClient(),
		Logger:           Logger(),
		Handlers:         Handlers(),
	}
}

// HTTPClient will return a new HTTP Client configured for the SDK.
//
// Does not use http.DefaultClient nor http.DefaultTransport.
func HTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{},
	}
}

// Handlers returns the default request handlers.
//
// Generally you shouldn't need to use this method directly, but
// is available if you need to reset the request handlers of an
// existing service client or session.
func Handlers() aws.Handlers {
	var handlers aws.Handlers

	handlers.Validate.PushBackNamed(ValidateEndpointHandler)
	handlers.Validate.AfterEachFn = aws.HandlerListStopOnError
	handlers.Build.PushBackNamed(SDKVersionUserAgentHandler)
	handlers.Build.AfterEachFn = aws.HandlerListStopOnError
	handlers.Sign.PushBackNamed(BuildContentLengthHandler)
	handlers.Send.PushBackNamed(ValidateReqSigHandler)
	handlers.Send.PushBackNamed(SendHandler)
	handlers.AfterRetry.PushBackNamed(AfterRetryHandler)
	handlers.ValidateResponse.PushBackNamed(ValidateResponseHandler)

	return handlers
}

//// CredChain returns the default credential chain.
////
//// Generally you shouldn't need to use this method directly, but
//// is available if you need to reset the credentials of an
//// existing service client or session's Config.
//func CredChain(cfg *aws.Config, handlers aws.Handlers) *aws.Credentials {
//	return aws.NewCredentials(&aws.ChainProvider{
//		VerboseErrors: aws.BoolValue(cfg.CredentialsChainVerboseErrors),
//		Providers: []aws.Provider{
//			&aws.EnvProvider{},
//			&aws.SharedCredentialsProvider{Filename: "", Profile: ""},
//			RemoteCredProvider(*cfg, handlers),
//		},
//	})
//}
//
//const (
//	httpProviderEnvVar     = "AWS_CONTAINER_CREDENTIALS_FULL_URI"
//	ecsCredsProviderEnvVar = "AWS_CONTAINER_CREDENTIALS_RELATIVE_URI"
//)
//
//// RemoteCredProvider returns a credentials provider for the default remote
//// endpoints such as EC2 or ECS Roles.
//func RemoteCredProvider(cfg aws.Config, handlers aws.Handlers) aws.Provider {
//	if u := os.Getenv(httpProviderEnvVar); len(u) > 0 {
//		return localHTTPCredProvider(cfg, handlers, u)
//	}
//
//	if uri := os.Getenv(ecsCredsProviderEnvVar); len(uri) > 0 {
//		u := fmt.Sprintf("http://169.254.170.2%s", uri)
//		return httpCredProvider(cfg, handlers, u)
//	}
//
//	return ec2RoleProvider(cfg, handlers)
//}
//
//func localHTTPCredProvider(cfg aws.Config, handlers aws.Handlers, u string) aws.Provider {
//	var errMsg string
//
//	parsed, err := url.Parse(u)
//	if err != nil {
//		errMsg = fmt.Sprintf("invalid URL, %v", err)
//	} else if host := aws.URLHostname(parsed); !(host == "localhost" || host == "127.0.0.1") {
//		errMsg = fmt.Sprintf("invalid host address, %q, only localhost and 127.0.0.1 are valid.", host)
//	}
//
//	if len(errMsg) > 0 {
//		if cfg.Logger != nil {
//			cfg.Logger.Log("Ignoring, HTTP credential provider", errMsg, err)
//		}
//		return aws.ErrorProvider{
//			Err:          awserr.New("CredentialsEndpointError", errMsg, err),
//			ProviderName: endpointcreds.ProviderName,
//		}
//	}
//
//	return httpCredProvider(cfg, handlers, u)
//}
//
//func httpCredProvider(cfg aws.Config, handlers aws.Handlers, u string) aws.Provider {
//	return endpointcreds.NewProviderClient(cfg, handlers, u,
//		func(p *endpointcreds.Provider) {
//			p.ExpiryWindow = 5 * time.Minute
//		},
//	)
//}
//
//func ec2RoleProvider(cfg aws.Config, handlers aws.Handlers) aws.Provider {
//	resolver := cfg.EndpointResolver
//	if resolver == nil {
//		resolver = endpoints.DefaultResolver()
//	}
//
//	e, _ := resolver.EndpointFor(endpoints.Ec2metadataServiceID, "")
//	return &ec2rolecreds.EC2RoleProvider{
//		Client:       ec2metadata.NewClient(cfg, handlers, e.URL, e.SigningRegion),
//		ExpiryWindow: 5 * time.Minute,
//	}
//}
