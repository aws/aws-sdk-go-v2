package customizations

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/awslabs/smithy-go"
	"github.com/awslabs/smithy-go/middleware"
	"github.com/awslabs/smithy-go/transport/http"
)

// UpdateEndpointOptions provides the options for the UpdateEndpoint middleware setup.
type UpdateEndpointOptions struct {
	UsePathStyle  bool
	UseAccelerate bool

	Region string
}

// UpdateEndpoint adds the middleware to the middleware stack based on the UpdateEndpointOptions.
func UpdateEndpoint(stack *middleware.Stack, options UpdateEndpointOptions) {
	stack.Serialize.Add(&updateEndpointMiddleware{
		usePathStyle:  options.UsePathStyle,
		useAccelerate: options.UseAccelerate,
		region:        options.Region,
	}, middleware.Before)
}

type updateEndpointMiddleware struct {
	usePathStyle  bool
	useAccelerate bool

	region string
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
	if iface, ok := in.Parameters.(bucketGetter); ok {
		bucket := iface.GetBucket()
		if len(bucket) != 0 {
			// if bucket len is zero, we ignore the following customizations
			if u.useAccelerate {
				if u.usePathStyle {
					// TODO: log that accelerate is not compatible with aws.Config.S3ForcePathStyle, ignoring S3ForcePathStyle.
				}

				if !hostCompatibleBucketName(req.URL, bucket) {
					return out, metadata, &smithy.SerializationError{
						Err: fmt.Errorf("bucket name %s is not compatible with S3 Accelerate", bucket),
					}
				}

				parts := strings.Split(req.URL.Host, ".")
				if len(parts) < 3 {
					return out, metadata, &smithy.SerializationError{
						Err: fmt.Errorf("unable to update endpoint host for S3 accelerate, hostname invalid, %s",
							req.URL.Host),
					}
				}

				if parts[0] == "s3" || strings.HasPrefix(parts[0], "s3-") {
					parts[0] = "s3-accelerate"
				}
				for i := 1; i+1 < len(parts); i++ {
					if parts[i] == u.region {
						parts = append(parts[:i], parts[i+1:]...)
						break
					}
				}
				req.URL.Host = strings.Join(parts, ".")
				moveBucketNameToHost(req.URL, bucket)
			} else if !u.usePathStyle {
				if !hostCompatibleBucketName(req.URL, bucket) {
					// bucket name must be valid to put into the host
					return
				}
				moveBucketNameToHost(req.URL, bucket)
			}
		}
	}

	return next.HandleSerialize(ctx, in)
}

// updates endpoint to use virtual host styling
func moveBucketNameToHost(u *url.URL, bucket string) {
	u.Host = bucket + "." + u.Host
	removeBucketFromPath(u)
}

// remove bucket from url
func removeBucketFromPath(u *url.URL) {
	u.Path = strings.Replace(u.Path, "/{Bucket}", "", -1)
	if u.Path == "" {
		u.Path = "/"
	}
}

// bucketGetter is an accessor interface to grab the "Bucket" field from
// an S3 type.
type bucketGetter interface {
	GetBucket() string
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
	return reDomain.MatchString(bucket) &&
		!reIPAddress.MatchString(bucket) &&
		!strings.Contains(bucket, "..")
}
