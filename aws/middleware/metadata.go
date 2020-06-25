package middleware

import (
	"context"

	"github.com/awslabs/smithy-go/middleware"
)

// RegisterServiceMetadata registers metadata about the service and operation into the middleware context
// so that it is available at runtime for other middleware to introspect.
type RegisterServiceMetadata struct {
	ServiceName    string
	ServiceID      string
	EndpointPrefix string
	SigningName    string
	Region         string
	OperationName  string
}

// ID returns the middleware identifier.
func (s RegisterServiceMetadata) ID() string {
	return "RegisterServiceMetadata"
}

// HandleInitialize registers service metadata information into the middleware context, allowing for introspection.
func (s RegisterServiceMetadata) HandleInitialize(
	ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
) (out middleware.InitializeOutput, metadata middleware.Metadata, err error) {
	if len(s.ServiceName) > 0 {
		ctx = setServiceName(ctx, s.ServiceName)
	}
	if len(s.ServiceID) > 0 {
		ctx = setServiceID(ctx, s.ServiceID)
	}
	if len(s.EndpointPrefix) > 0 {
		ctx = setEndpointPrefix(ctx, s.EndpointPrefix)
	}
	if len(s.SigningName) > 0 {
		ctx = SetSigningName(ctx, s.SigningName)
	}
	if len(s.Region) > 0 {
		ctx = setRegion(ctx, s.Region)
	}
	if len(s.OperationName) > 0 {
		ctx = setOperationName(ctx, s.OperationName)
	}
	return next.HandleInitialize(ctx, in)
}

// service metadata keys for storing and lookup of runtime stack information.
type (
	serviceNameKey    struct{}
	serviceIDKey      struct{}
	endpointPrefixKey struct{}
	signingNameKey    struct{}
	signingRegionKey  struct{}
	regionKey         struct{}
	operationNameKey  struct{}
)

// GetServiceName retrieves the service name from the context.
func GetServiceName(ctx context.Context) (v string) {
	v, _ = ctx.Value(serviceNameKey{}).(string)
	return v
}

// GetServiceID retrieves the service id from the context.
func GetServiceID(ctx context.Context) (v string) {
	v, _ = ctx.Value(serviceIDKey{}).(string)
	return v
}

// GetEndpointPrefix retrieves the service endpoints id from the context.
func GetEndpointPrefix(ctx context.Context) (v string) {
	v, _ = ctx.Value(endpointPrefixKey{}).(string)
	return v
}

// GetSigningName retrieves the service signing name from the context.
func GetSigningName(ctx context.Context) (v string) {
	v, _ = ctx.Value(signingNameKey{}).(string)
	return v
}

// GetSigningRegion retrieves the region from the context.
func GetSigningRegion(ctx context.Context) (v string) {
	v, _ = ctx.Value(signingRegionKey{}).(string)
	return v
}

// GetRegion retrieves the endpoint region from the context.
func GetRegion(ctx context.Context) (v string) {
	v, _ = ctx.Value(regionKey{}).(string)
	return v
}

// GetOperationName retrieves the service operation metadata from the context.
func GetOperationName(ctx context.Context) (v string) {
	v, _ = ctx.Value(operationNameKey{}).(string)
	return v
}

// SetSigningName set or modifies the signing name on the context.
func SetSigningName(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, signingNameKey{}, value)
}

// SetSigningRegion sets or modifies the region on the context.
func SetSigningRegion(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, signingRegionKey{}, value)
}

// setServiceName sets the service name on the context.
func setServiceName(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, serviceNameKey{}, value)
}

// setServiceID sets the service id on the context.
func setServiceID(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, serviceIDKey{}, value)
}

// setEndpointPrefix sets the service endpoint id on the context.
func setEndpointPrefix(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, endpointPrefixKey{}, value)
}

// setRegion sets the endpoint region on the context.
func setRegion(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, regionKey{}, value)
}

// setOperationName sets the service operation on the context.
func setOperationName(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, operationNameKey{}, value)
}
