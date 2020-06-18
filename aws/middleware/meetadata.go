package middleware

import (
	"context"

	"github.com/awslabs/smithy-go/middleware"
)

// RegisterServiceMetadata registers metadata about the service and operation into the middleware context
// so that it is available at runtime for other middleware to introspect.
type RegisterServiceMetadata struct {
	ServiceName string
	ServiceID   string
	EndpointsID string
	SigningName string
	Region      string
	Operation   OperationMetadata
}

// OperationMetadata metadata about the service operation.
type OperationMetadata struct {
	Name     string
	HTTPPath string
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
	if len(s.EndpointsID) > 0 {
		ctx = setEndpointID(ctx, s.EndpointsID)
	}
	if len(s.SigningName) > 0 {
		ctx = SetSigningName(ctx, s.SigningName)
	}
	if len(s.Region) > 0 {
		ctx = setRegion(ctx, s.Region)
	}
	if s.Operation != (OperationMetadata{}) {
		ctx = setOperationMetadata(ctx, s.Operation)
	}
	return next.HandleInitialize(ctx, in)
}

// service metadata keys for storing and lookup of runtime stack information.
type (
	serviceNameKey       struct{}
	serviceIDKey         struct{}
	endpointIDKey        struct{}
	signingNameKey       struct{}
	signingRegionKey     struct{}
	regionKey            struct{}
	operationMetadataKey struct{}
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

// GetEndpointID retrieves the service endpoints id from the context.
func GetEndpointID(ctx context.Context) (v string) {
	v, _ = ctx.Value(endpointIDKey{}).(string)
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

// GetOperationMetadata retrieves the service operation metadata from the context.
func GetOperationMetadata(ctx context.Context) (v OperationMetadata) {
	v, _ = ctx.Value(operationMetadataKey{}).(OperationMetadata)
	return v
}

// SetSigningName set or modifies the signing name on the context.
func SetSigningName(ctx context.Context, value string) context.Context {
	ctx = context.WithValue(ctx, signingNameKey{}, value)
	return ctx
}

// SetSigningRegion sets or modifies the region on the context.
func SetSigningRegion(ctx context.Context, value string) context.Context {
	ctx = context.WithValue(ctx, signingRegionKey{}, value)
	return ctx
}

// setServiceName sets the service name on the context.
func setServiceName(ctx context.Context, value string) context.Context {
	ctx = context.WithValue(ctx, serviceNameKey{}, value)
	return ctx
}

// setServiceID sets the service id on the context.
func setServiceID(ctx context.Context, value string) context.Context {
	ctx = context.WithValue(ctx, serviceIDKey{}, value)
	return ctx
}

// setEndpointID sets the service endpoint id on the context.
func setEndpointID(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, endpointIDKey{}, value)
}

// setRegion sets the endpoint region on the context.
func setRegion(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, regionKey{}, value)
}

// setOperationMetadata sets the service operation on the context.
func setOperationMetadata(ctx context.Context, value OperationMetadata) context.Context {
	return context.WithValue(ctx, operationMetadataKey{}, value)
}
