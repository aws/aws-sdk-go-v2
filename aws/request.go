package aws

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/awserr"
)

const (
	// ErrCodeSerialization is the serialization error code that is received
	// during protocol unmarshaling.
	ErrCodeSerialization = "SerializationError"

	// ErrCodeRead is an error that is returned during HTTP reads.
	ErrCodeRead = "ReadError"

	// ErrCodeResponseTimeout is the connection timeout error that is received
	// during body reads.
	ErrCodeResponseTimeout = "ResponseTimeout"

	// CanceledErrorCode is the error code that will be returned by an
	// API request that was canceled. Requests given a Context may
	// return this error when canceled.
	CanceledErrorCode = "RequestCanceled"
)

// A Request is the service request to be made.
type Request struct {
	Config   Config
	Metadata Metadata
	Handlers Handlers

	Retryer
	Time                   time.Time
	ExpireTime             time.Duration
	Operation              *Operation
	HTTPRequest            *http.Request
	HTTPResponse           *http.Response
	Body                   io.ReadSeeker
	BodyStart              int64 // offset from beginning of Body that the request body starts
	Params                 interface{}
	Error                  error
	Data                   interface{}
	RequestID              string
	RetryCount             int
	Retryable              *bool
	RetryDelay             time.Duration
	NotHoist               bool
	SignedHeaderVals       http.Header
	LastSignedAt           time.Time
	DisableFollowRedirects bool

	context Context

	built bool

	// Need to persist an intermediate body between the input Body and HTTP
	// request body because the HTTP Client's transport can maintain a reference
	// to the HTTP request's body after the client has returned. This value is
	// safe to use concurrently and wrap the input Body for each HTTP request.
	safeBody *offsetReader
}

// An Operation is the service API operation to be made.
type Operation struct {
	Name       string
	HTTPMethod string
	HTTPPath   string
	*Paginator

	BeforePresignFn func(r *Request) error
}

// New returns a new Request pointer for the service API
// operation and parameters.
//
// Params is any value of input parameters to be the request payload.
// Data is pointer value to an object which the request's response
// payload will be deserialized to.
func New(cfg Config, metadata Metadata, handlers Handlers,
	retryer Retryer, operation *Operation, params interface{}, data interface{}) *Request {

	// TODO improve this experiance for config copy?
	cfg = cfg.Copy()

	method := operation.HTTPMethod
	if method == "" {
		method = "POST"
	}

	httpReq, _ := http.NewRequest(method, "", nil)

	// TODO need better way of handling this error... NewRequest should return error.
	endpoint, err := cfg.EndpointResolver.ResolveEndpoint(metadata.ServiceName, cfg.Region)
	if err == nil {
		// TODO so ugly
		metadata.Endpoint = endpoint.URL
		if len(endpoint.SigningName) > 0 {
			metadata.SigningName = endpoint.SigningName
		}
		if len(endpoint.SigningRegion) > 0 {
			metadata.SigningRegion = endpoint.SigningRegion
		}

		httpReq.URL, err = url.Parse(endpoint.URL + operation.HTTPPath)
		if err != nil {
			httpReq.URL = &url.URL{}
			err = awserr.New("InvalidEndpointURL", "invalid endpoint uri", err)
		}
	}

	r := &Request{
		Config:   cfg,
		Metadata: metadata,
		Handlers: handlers.Copy(),

		Retryer:     retryer,
		Time:        time.Now(),
		ExpireTime:  0,
		Operation:   operation,
		HTTPRequest: httpReq,
		Body:        nil,
		Params:      params,
		Error:       err,
		Data:        data,
	}
	r.SetBufferBody([]byte{})

	return r
}

// A Option is a functional option that can augment or modify a request when
// using a WithContext API operation method.
type Option func(*Request)

// WithGetResponseHeader builds a request Option which will retrieve a single
// header value from the HTTP Response. If there are multiple values for the
// header key use WithGetResponseHeaders instead to access the http.Header
// map directly. The passed in val pointer must be non-nil.
//
// This Option can be used multiple times with a single API operation.
//
//    var id2, versionID string
//    svc.PutObjectWithContext(ctx, params,
//        request.WithGetResponseHeader("x-amz-id-2", &id2),
//        request.WithGetResponseHeader("x-amz-version-id", &versionID),
//    )
func WithGetResponseHeader(key string, val *string) Option {
	return func(r *Request) {
		r.Handlers.Complete.PushBack(func(req *Request) {
			*val = req.HTTPResponse.Header.Get(key)
		})
	}
}

// WithGetResponseHeaders builds a request Option which will retrieve the
// headers from the HTTP response and assign them to the passed in headers
// variable. The passed in headers pointer must be non-nil.
//
//    var headers http.Header
//    svc.PutObjectWithContext(ctx, params, request.WithGetResponseHeaders(&headers))
func WithGetResponseHeaders(headers *http.Header) Option {
	return func(r *Request) {
		r.Handlers.Complete.PushBack(func(req *Request) {
			*headers = req.HTTPResponse.Header
		})
	}
}

// WithLogLevel is a request option that will set the request to use a specific
// log level when the request is made.
//
//     svc.PutObjectWithContext(ctx, params, request.WithLogLevel(LogDebugWithHTTPBody)
func WithLogLevel(l LogLevel) Option {
	return func(r *Request) {
		r.Config.LogLevel = l
	}
}

// ApplyOptions will apply each option to the request calling them in the order
// the were provided.
func (r *Request) ApplyOptions(opts ...Option) {
	for _, opt := range opts {
		opt(r)
	}
}

// Context will always returns a non-nil context. If Request does not have a
// context BackgroundContext will be returned.
func (r *Request) Context() Context {
	if r.context != nil {
		return r.context
	}
	return BackgroundContext()
}

// SetContext adds a Context to the current request that can be used to cancel
// a in-flight request. The Context value must not be nil, or this method will
// panic.
//
// Unlike http.Request.WithContext, SetContext does not return a copy of the
// Request. It is not safe to use use a single Request value for multiple
// requests. A new Request should be created for each API operation request.
//
// Go 1.6 and below:
// The http.Request's Cancel field will be set to the Done() value of
// the context. This will overwrite the Cancel field's value.
//
// Go 1.7 and above:
// The http.Request.WithContext will be used to set the context on the underlying
// http.Request. This will create a shallow copy of the http.Request. The SDK
// may create sub contexts in the future for nested requests such as retries.
func (r *Request) SetContext(ctx Context) {
	if ctx == nil {
		panic("context cannot be nil")
	}
	setRequestContext(r, ctx)
}

// WillRetry returns if the request's can be retried.
func (r *Request) WillRetry() bool {
	return r.Error != nil && BoolValue(r.Retryable) && r.RetryCount < r.MaxRetries()
}

// ParamsFilled returns if the request's parameters have been populated
// and the parameters are valid. False is returned if no parameters are
// provided or invalid.
func (r *Request) ParamsFilled() bool {
	return r.Params != nil && reflect.ValueOf(r.Params).Elem().IsValid()
}

// setbufferbody will set the request's body bytes that will be sent to
// the service api.
func (r *request) setbufferbody(buf []byte) {
	r.setreaderbody(bytes.newreader(buf))
}

// setstringbody sets the body of the request to be backed by a string.
func (r *request) setstringbody(s string) {
	r.setreaderbody(strings.newreader(s))
}

// setreaderbody will set the request's body reader.
func (r *request) setreaderbody(reader io.readseeker) {
	r.body = reader
	r.resetbody()
}

// presign returns the request's signed url. error will be returned
// if the signing fails.
func (r *request) presign(expiretime time.duration) (string, error) {
	r.expiretime = expiretime
	r.nothoist = false

	if r.operation.beforepresignfn != nil {
		r = r.copy()
		err := r.operation.beforepresignfn(r)
		if err != nil {
			return "", err
		}
	}

	r.sign()
	if r.error != nil {
		return "", r.error
	}
	return r.httprequest.url.string(), nil
}

// presignrequest behaves just like presign, with the addition of returning a
// set of headers that were signed.
//
// returns the url string for the api operation with signature in the query string,
// and the http headers that were included in the signature. these headers must
// be included in any http request made with the presigned url.
//
// to prevent hoisting any headers to the query string set nothoist to true on
// this request value prior to calling presignrequest.
func (r *request) presignrequest(expiretime time.duration) (string, http.header, error) {
	r.expiretime = expiretime
	r.sign()
	if r.error != nil {
		return "", nil, r.error
	}
	return r.httprequest.url.string(), r.signedheadervals, nil
}

func debuglogreqerror(r *request, stage string, retrying bool, err error) {
	if !r.config.loglevel.matches(logdebugwithrequesterrors) {
		return
	}

	retrystr := "not retrying"
	if retrying {
		retrystr = "will retry"
	}

	r.config.logger.log(fmt.sprintf("debug: %s %s/%s failed, %s, error %v",
		stage, r.metadata.servicename, r.operation.name, retrystr, err))
}

// build will build the request's object so it can be signed and sent
// to the service. build will also validate all the request's parameters.
// anny additional build handlers set on this request will be run
// in the order they were set.
//
// the request will only be built once. multiple calls to build will have
// no effect.
//
// if any validate or build errors occur the build will stop and the error
// which occurred will be returned.
func (r *request) build() error {
	if !r.built {
		r.handlers.validate.run(r)
		if r.error != nil {
			debuglogreqerror(r, "validate request", false, r.error)
			return r.error
		}
		r.handlers.build.run(r)
		if r.error != nil {
			debuglogreqerror(r, "build request", false, r.error)
			return r.error
		}
		r.built = true
	}

	return r.error
}

// sign will sign the request returning error if errors are encountered.
//
// send will build the request prior to signing. all sign handlers will
// be executed in the order they were set.
func (r *request) sign() error {
	r.build()
	if r.error != nil {
		debuglogreqerror(r, "build request", false, r.error)
		return r.error
	}

	r.handlers.sign.run(r)
	return r.error
}

func (r *request) getnextrequestbody() (io.readcloser, error) {
	if r.safebody != nil {
		r.safebody.close()
	}

	r.safebody = newoffsetreader(r.body, r.bodystart)

	// go 1.8 tightened and clarified the rules code needs to use when building
	// requests with the http package. go 1.8 removed the automatic detection
	// of if the request.body was empty, or actually had bytes in it. the sdk
	// always sets the request.body even if it is empty and should not actually
	// be sent. this is incorrect.
	//
	// go 1.8 did add a http.nobody value that the sdk can use to tell the http
	// client that the request really should be sent without a body. the
	// request.body cannot be set to nil, which is preferable, because the
	// field is exported and could introduce nil pointer dereferences for users
	// of the sdk if they used that field.
	//
	// related golang/go#18257
	l, err := computebodylength(r.body)
	if err != nil {
		return nil, awserr.new(errcodeserialization, "failed to compute request body size", err)
	}

	var body io.readcloser
	if l == 0 {
		body = nobody
	} else if l > 0 {
		body = r.safebody
	} else {
		// hack to prevent sending bodies for methods where the body
		// should be ignored by the server. sending bodies on these
		// methods without an associated contentlength will cause the
		// request to socket timeout because the server does not handle
		// transfer-encoding: chunked bodies for these methods.
		//
		// this would only happen if a readerseekercloser was used with
		// a io.reader that was not also an io.seeker.
		switch r.operation.httpmethod {
		case "get", "head", "delete":
			body = nobody
		default:
			body = r.safebody
		}
	}

	return body, nil
}

// attempts to compute the length of the body of the reader using the
// io.seeker interface. if the value is not seekable because of being
// a readerseekercloser without an unerlying seeker -1 will be returned.
// if no error occurs the length of the body will be returned.
func computebodylength(r io.readseeker) (int64, error) {
	seekable := true
	// determine if the seeker is actually seekable. readerseekercloser
	// hides the fact that a io.readers might not actually be seekable.
	switch v := r.(type) {
	case readerseekercloser:
		seekable = v.isseeker()
	case *readerseekercloser:
		seekable = v.isseeker()
	}
	if !seekable {
		return -1, nil
	}

	curoffset, err := r.seek(0, 1)
	if err != nil {
		return 0, err
	}

	endoffset, err := r.seek(0, 2)
	if err != nil {
		return 0, err
	}

	_, err = r.seek(curoffset, 0)
	if err != nil {
		return 0, err
	}

	return endoffset - curoffset, nil
}

// getbody will return an io.readseeker of the request's underlying
// input body with a concurrency safe wrapper.
func (r *request) getbody() io.readseeker {
	return r.safebody
}

// send will send the request returning error if errors are encountered.
//
// send will sign the request prior to sending. all send handlers will
// be executed in the order they were set.
//
// canceling a request is non-deterministic. if a request has been canceled,
// then the transport will choose, randomly, one of the state channels during
// reads or getting the connection.
//
// readloop() and getconn(req *request, cm connectmethod)
// https://github.com/golang/go/blob/master/src/net/http/transport.go
//
// send will not close the request.request's body.
func (r *request) send() error {
	defer func() {
		// regardless of success or failure of the request trigger the complete
		// request handlers.
		r.handlers.complete.run(r)
	}()

	for {
		if boolvalue(r.retryable) {
			if r.config.loglevel.matches(logdebugwithrequestretries) {
				r.config.logger.log(fmt.sprintf("debug: retrying request %s/%s, attempt %d",
					r.metadata.servicename, r.operation.name, r.retrycount))
			}

			// the previous http.request will have a reference to the r.body
			// and the http client's transport may still be reading from
			// the request's body even though the client's do returned.
			r.httprequest = copyhttprequest(r.httprequest, nil)
			r.resetbody()

			// closing response body to ensure that no response body is leaked
			// between retry attempts.
			if r.httpresponse != nil && r.httpresponse.body != nil {
				r.httpresponse.body.close()
			}
		}

		r.sign()
		if r.error != nil {
			return r.error
		}

		r.retryable = nil

		r.handlers.send.run(r)
		if r.error != nil {
			if !shouldretrycancel(r) {
				return r.error
			}

			err := r.error
			r.handlers.retry.run(r)
			r.handlers.afterretry.run(r)
			if r.error != nil {
				debuglogreqerror(r, "send request", false, err)
				return r.error
			}
			debuglogreqerror(r, "send request", true, err)
			continue
		}
		r.handlers.unmarshalmeta.run(r)
		r.handlers.validateresponse.run(r)
		if r.error != nil {
			r.handlers.unmarshalerror.run(r)
			err := r.error

			r.handlers.retry.run(r)
			r.handlers.afterretry.run(r)
			if r.error != nil {
				debuglogreqerror(r, "validate response", false, err)
				return r.error
			}
			debuglogreqerror(r, "validate response", true, err)
			continue
		}

		r.handlers.unmarshal.run(r)
		if r.error != nil {
			err := r.error
			r.handlers.retry.run(r)
			r.handlers.afterretry.run(r)
			if r.error != nil {
				debuglogreqerror(r, "unmarshal response", false, err)
				return r.error
			}
			debuglogreqerror(r, "unmarshal response", true, err)
			continue
		}

		break
	}

	return nil
}

// copy will copy a request which will allow for local manipulation of the
// request.
func (r *request) copy() *request {
	req := &request{}
	*req = *r
	req.handlers = r.handlers.copy()
	op := *r.operation
	req.operation = &op
	return req
}

// addtouseragent adds the string to the end of the request's current user agent.
func addtouseragent(r *request, s string) {
	curua := r.httprequest.header.get("user-agent")
	if len(curua) > 0 {
		s = curua + " " + s
	}
	r.httprequest.header.set("user-agent", s)
}

func shouldretrycancel(r *request) bool {
	awserr, ok := r.error.(awserr.error)
	timeouterr := false
	errstr := r.error.error()
	if ok {
		if awserr.code() == cancelederrorcode {
			return false
		}
		err := awserr.origerr()
		neterr, netok := err.(net.error)
		timeouterr = netok && neterr.temporary()
		if urlerr, ok := err.(*url.error); !timeouterr && ok {
			errstr = urlerr.err.error()
		}
	}

	// there can be two types of canceled errors here.
	// the first being a net.error and the other being an error.
	// if the request was timed out, we want to continue the retry
	// process. otherwise, return the canceled error.
	return timeouterr ||
		(errstr != "net/http: request canceled" &&
			errstr != "net/http: request canceled while waiting for connection")

}
