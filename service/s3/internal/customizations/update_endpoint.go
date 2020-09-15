package customizations

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/awslabs/smithy-go/middleware"
	"github.com/awslabs/smithy-go/transport/http"
)

// UpdateEndpointOptions provides the options for the UpdateEndpoint middleware setup.
type UpdateEndpointOptions struct {
	// functional pointer to fetch bucket name from provided input.
	// The function is intended to take an input value, and
	// return a string pointer to value of string, and bool if
	// input has no bucket member.
	GetBucketFromInput func(interface{}) (*string, bool)

	// use path style
	UsePathStyle bool
}

// UpdateEndpoint adds the middleware to the middleware stack based on the UpdateEndpointOptions.
func UpdateEndpoint(stack *middleware.Stack, options UpdateEndpointOptions) {
	stack.Serialize.Insert(&updateEndpointMiddleware{
		getBucketFromInput: options.GetBucketFromInput,
		usePathStyle:       options.UsePathStyle,
	}, "OperationSerializer", middleware.After)
}

type updateEndpointMiddleware struct {
	getBucketFromInput func(interface{}) (*string, bool)
	usePathStyle       bool
}

// ID returns the middleware ID.
func (*updateEndpointMiddleware) ID() string { return "S3:UpdateEndpointMiddleware" }

func (u *updateEndpointMiddleware) HandleSerialize(
	ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler,
) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	req, ok := in.Request.(*http.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown request type %T", req)
	}

	// Below customization only apply if bucket name is provided
	bucket, ok := u.getBucketFromInput(in.Parameters)
	if ok && bucket != nil {
		if err := u.updateEndpointFromConfig(req, *bucket); err != nil {
			return out, metadata, err
		}
	}
	return next.HandleSerialize(ctx, in)
}

func (u updateEndpointMiddleware) updateEndpointFromConfig(req *http.Request, bucket string) error {
	if !u.usePathStyle {
		if !hostCompatibleBucketName(req.URL, bucket) {
			// bucket name must be valid to put into the host
			return fmt.Errorf("bucket name %s is not compatible with S3", bucket)
		}
		moveBucketNameToHost(req.URL, bucket)
	}
	return nil
}

// updates endpoint to use virtual host styling
func moveBucketNameToHost(u *url.URL, bucket string) {
	u.Host = bucket + "." + u.Host
	removeBucketFromPath(u, bucket)
}

// remove bucket from url
func removeBucketFromPath(u *url.URL, bucket string) {
	u.Path = strings.Replace(u.Path, "/"+bucket, "", -1)
	if u.Path == "" {
		u.Path = "/"
	}
}

// hostCompatibleBucketName returns true if the request should
// put the bucket in the host. This is false if S3ForcePathStyle is
// explicitly set or if the bucket is not DNS compatible.
func hostCompatibleBucketName(u *url.URL, bucket string) bool {
	// Bucket might be DNS compatible but dots in the hostname will fail
	// certificate validation, so do not use host-style.
	if u.Scheme == "https" && strings.Contains(bucket, ".") {
		return false
	}

	// if the bucket is DNS compatible
	return dnsCompatibleBucketName(bucket)
}

var reDomain = regexp.MustCompile(`^[a-z0-9][a-z0-9\.\-]{1,61}[a-z0-9]$`)
var reIPAddress = regexp.MustCompile(`^(\d+\.){3}\d+$`)

// dnsCompatibleBucketName returns true if the bucket name is DNS compatible.
// Buckets created outside of the classic region MUST be DNS compatible.
func dnsCompatibleBucketName(bucket string) bool {
	if strings.Contains(bucket, "..") {
		return false
	}

	// checks for `^[a-z0-9][a-z0-9\.\-]{1,61}[a-z0-9]$` domain mapping
	if !((bucket[0] > 96 && bucket[0] < 123) || (bucket[0] > 47 && bucket[0] < 58)) {
		return false
	}

	for _, c := range bucket[1:] {
		if !((c > 96 && c < 123) || (c > 47 && c < 58) || c == 46 || c == 45) {
			return false
		}
	}

	// checks for `^(\d+\.){3}\d+$` IPaddressing
	v := strings.SplitN(bucket, ".", -1)
	if len(v) == 4 {
		for _, c := range bucket {
			if !((c > 47 && c < 58) || c == 46) {
				// we confirm that this is not a IP address
				return true
			}
		}
		// this is a IP address
		return false
	}

	return true
}
