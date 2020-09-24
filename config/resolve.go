package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// ResolveDefaultAWSConfig will write default configuration values into the cfg
// value. It will write the default values, overwriting any previous value.
//
// This should be used as the first resolver in the slice of resolvers when
// resolving external configuration.
func ResolveDefaultAWSConfig(cfg *aws.Config, configs Configs) error {
	*cfg = aws.Config{
		Credentials: aws.AnonymousCredentials{},
	}
	return nil
}

// ResolveCustomCABundle extracts the first instance of a custom CA bundle filename
// from the external configurations. It will update the HTTP Client's builder
// to be configured with the custom CA bundle.
//
// Config provider used:
// * CustomCABundleProvider
func ResolveCustomCABundle(cfg *aws.Config, configs Configs) error {
	pemCerts, found, err := GetCustomCABundle(configs)
	if err != nil {
		// TODO error handling, What is the best way to handle this?
		// capture previous errors continue. error out if all errors
		return err
	}
	if !found {
		return nil
	}

	type withTransportOptions interface {
		WithTransportOptions(...func(*http.Transport)) aws.HTTPClient
	}

	trOpts, ok := cfg.HTTPClient.(withTransportOptions)
	if !ok {
		return fmt.Errorf("unable to add custom RootCAs HTTPClient, "+
			"has no WithTransportOptions, %T", cfg.HTTPClient)
	}

	var appendErr error
	client := trOpts.WithTransportOptions(func(tr *http.Transport) {
		if tr.TLSClientConfig == nil {
			tr.TLSClientConfig = &tls.Config{}
		}
		if tr.TLSClientConfig.RootCAs == nil {
			tr.TLSClientConfig.RootCAs = x509.NewCertPool()
		}
		if !tr.TLSClientConfig.RootCAs.AppendCertsFromPEM(pemCerts) {
			appendErr = fmt.Errorf("failed to load custom CA bundle PEM file")
		}
	})
	if appendErr != nil {
		return appendErr
	}

	cfg.HTTPClient = client
	return err
}

// ResolveRegion extracts the first instance of a Region from the Configs slice.
//
// Config providers used:
// * RegionProvider
func ResolveRegion(cfg *aws.Config, configs Configs) error {
	v, found, err := GetRegion(configs)
	if err != nil {
		// TODO error handling, What is the best way to handle this?
		// capture previous errors continue. error out if all errors
		return err
	}
	if !found {
		return nil
	}

	cfg.Region = v
	return nil
}

// ResolveDefaultRegion extracts the first instance of a default region and sets `aws.Config.Region` to the default
// region if region had not been resolved from other sources.
func ResolveDefaultRegion(cfg *aws.Config, configs Configs) error {
	if len(cfg.Region) > 0 {
		return nil
	}

	region, found, err := GetDefaultRegion(configs)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}

	cfg.Region = region

	return nil
}

// ResolveHTTPClient extracts the first instance of a HTTPClient and sets `aws.Config.HTTPClient` to the HTTPClient instance
// if one has not been resolved from other sources.
func ResolveHTTPClient(cfg *aws.Config, configs Configs) error {
	c, found, err := GetHTTPClient(configs)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}

	cfg.HTTPClient = c
	return nil
}

// ResolveAPIOptions extracts the first instance of APIOptions and sets `aws.Config.APIOptions` to the resolved API options
// if one has not been resolved from other sources.
func ResolveAPIOptions(cfg *aws.Config, configs Configs) error {
	o, found, err := GetAPIOptions(configs)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}

	cfg.APIOptions = o

	return nil
}

// ResolveEndpointResolver extracts the first instance of a EndpointResolverFunc from the config slice
// and sets the functions result on the aws.Config.EndpointResolver
func ResolveEndpointResolver(cfg *aws.Config, configs Configs) error {
	endpointResolver, found, err := GetEndpointResolver(configs)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}

	cfg.EndpointResolver = endpointResolver

	return nil
}
