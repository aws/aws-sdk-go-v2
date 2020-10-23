package customizations

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/awslabs/smithy-go/middleware"
	smithyid "github.com/awslabs/smithy-go/middleware/id"
	"github.com/awslabs/smithy-go/transport/http"

	"github.com/aws/aws-sdk-go-v2/service/internal/s3shared"
)

// UpdateEndpointOptions provides the options for the UpdateEndpoint middleware setup.
type UpdateEndpointOptions struct {
	// region used
	Region string

	// functional pointer to fetch bucket name from provided input.
	// The function is intended to take an input value, and
	// return a string pointer to value of string, and bool if
	// input has no bucket member.
	GetBucketFromInput func(interface{}) (*string, bool)

	// use path style
	UsePathStyle bool

	// use transfer acceleration
	UseAccelerate bool

	// functional pointer to indicate support for accelerate.
	// The function is intended to take an input value, and
	// return if the operation supports accelerate.
	SupportsAccelerate func(interface{}) bool

	// use dualstack
	UseDualstack bool
}

// UpdateEndpoint adds the middleware to the middleware stack based on the UpdateEndpointOptions.
func UpdateEndpoint(stack *middleware.Stack, options UpdateEndpointOptions) error {
	// enable dual stack support
	if err := stack.Serialize.Insert(&s3shared.EnableDualstackMiddleware{
		UseDualstack: options.UseDualstack,
		ServiceID:    "s3",
	}, smithyid.OperationSerializer, middleware.After); err != nil {
		return err
	}

	// update endpoint to use options for path style and accelerate
	return stack.Serialize.Insert(&updateEndpointMiddleware{
		region:             options.Region,
		usePathStyle:       options.UsePathStyle,
		getBucketFromInput: options.GetBucketFromInput,
		useAccelerate:      options.UseAccelerate,
		supportsAccelerate: options.SupportsAccelerate,
	}, (&s3shared.EnableDualstackMiddleware{}).ID(), middleware.After)
}

type updateEndpoint struct {
	region string

	// path style options
	usePathStyle       bool
	getBucketFromInput func(interface{}) (*string, bool)

	// accelerate options
	useAccelerate      bool
	supportsAccelerate func(interface{}) bool
}

// ID returns the middleware ID.
func (*updateEndpoint) ID() string { return "S3:UpdateEndpoint" }

func (u *updateEndpoint) HandleSerialize(
	ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler,
) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	req, ok := in.Request.(*http.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown request type %T", req)
	}

	// check if accelerate is supported
	if u.supportsAccelerate == nil || (u.useAccelerate && !u.supportsAccelerate(in.Parameters)) {
		// accelerate is not supported, thus will be ignored
		log.Println("Transfer acceleration is not supported for the operation, ignoring UseAccelerate.")
		u.useAccelerate = false
	}

	// transfer acceleration is not supported with path style urls
	if u.useAccelerate && u.usePathStyle {
		log.Println("UseAccelerate is not compatible with UsePathStyle, ignoring UsePathStyle.")
		u.usePathStyle = false
	}

	if u.getBucketFromInput != nil {
		// Below customization only apply if bucket name is provided
		bucket, ok := u.getBucketFromInput(in.Parameters)
		if ok && bucket != nil {
			if err := u.updateEndpointFromConfig(req, *bucket); err != nil {
				return out, metadata, err
			}
		}
	}

	return next.HandleSerialize(ctx, in)
}

func (u updateEndpoint) updateEndpointFromConfig(req *http.Request, bucket string) error {
	// do nothing if path style is enforced
	if u.usePathStyle {
		return nil
	}

	if !hostCompatibleBucketName(req.URL, bucket) {
		// bucket name must be valid to put into the host
		return fmt.Errorf("bucket name %s is not compatible with S3", bucket)
	}

	// accelerate is only supported if use path style is disabled
	if u.useAccelerate {
		parts := strings.Split(req.URL.Host, ".")
		if len(parts) < 3 {
			return fmt.Errorf("unable to update endpoint host for S3 accelerate, hostname invalid, %s", req.URL.Host)
		}

		if parts[0] == "s3" || strings.HasPrefix(parts[0], "s3-") {
			parts[0] = "s3-accelerate"
		}

		for i := 1; i+1 < len(parts); i++ {
			if strings.EqualFold(parts[i], u.region) {
				parts = append(parts[:i], parts[i+1:]...)
				break
			}
		}

		// construct the url host
		req.URL.Host = strings.Join(parts, ".")
	}

	// move bucket to follow virtual host style
	moveBucketNameToHost(req.URL, bucket)
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
