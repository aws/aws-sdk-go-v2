package client

import (
	"github.com/awslabs/smithy-go/middleware"
	"github.com/awslabs/smithy-go/transport/http"
)

var _ http.Header
var _ middleware.MiddlewareHandler
