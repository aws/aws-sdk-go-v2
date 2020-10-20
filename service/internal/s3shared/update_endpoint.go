package s3shared

import (
	"context"
	"fmt"
	"strings"

	"github.com/awslabs/smithy-go/middleware"
	"github.com/awslabs/smithy-go/transport/http"
)

type EnableDualstackMiddleware struct {
	UseDualstack bool
	ServiceID    string
}

// ID returns the middleware ID.
func (*EnableDualstackMiddleware) ID() string { return "EnableDualstackMiddleware" }

func (u *EnableDualstackMiddleware) HandleSerialize(
	ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler,
) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	req, ok := in.Request.(*http.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown request type %T", req)
	}

	if u.UseDualstack {
		parts := strings.Split(req.URL.Host, ".")
		if len(parts) < 3 {
			return out, metadata, fmt.Errorf("unable to update endpoint host for dualstack, hostname invalid, %s", req.URL.Host)
		}

		for i := 0; i+1 < len(parts); i++ {
			if strings.EqualFold(parts[i], u.ServiceID) {
				parts[i] = parts[i] + ".dualstack"
				break
			}
		}

		// construct the url host
		req.URL.Host = strings.Join(parts, ".")
	}

	return next.HandleSerialize(ctx, in)
}
