package aws

// A ClientConfig provides configuration to a service client instance.
type ClientConfig struct {
	Config        *Config
	Handlers      Handlers
	Endpoint      string
	SigningRegion string
	SigningName   string
}

// ConfigProvider provides a generic way for a service client to receive
// the ClientConfig without circular dependencies.
type ConfigProvider interface {
	ClientConfig(serviceName string, cfgs ...*Config) ClientConfig
}

// A Client implements the base client request and response handling
// used by all service clients.
type Client struct {
	Retryer
	ClientInfo

	Config   Config
	Handlers Handlers
}

// NewClient will return a pointer to a new initialized service client.
func NewClient(cfg Config, info ClientInfo, handlers Handlers, options ...func(*Client)) *Client {
	svc := &Client{
		Config:     cfg,
		ClientInfo: info,
		Handlers:   handlers.Copy(),
	}

	retryer := cfg.Retryer
	if retryer == nil {
		// TODO need better way of specifing default num retries
		retryer = DefaultRetryer{NumMaxRetries: 3}
	}
	svc.Retryer = retryer

	svc.AddDebugHandlers()

	for _, option := range options {
		option(svc)
	}

	return svc
}

// NewRequest returns a new Request pointer for the service API
// operation and parameters.
func (c *Client) NewRequest(operation *Operation, params interface{}, data interface{}) *Request {
	return New(c.Config, c.ClientInfo, c.Handlers, c.Retryer, operation, params, data)
}

// AddDebugHandlers injects debug logging handlers into the service to log request
// debug information.
func (c *Client) AddDebugHandlers() {
	if !c.Config.LogLevel.AtLeast(LogDebug) {
		return
	}

	c.Handlers.Send.PushFrontNamed(NamedHandler{Name: "awssdk.client.LogRequest", Fn: logRequest})
	c.Handlers.Send.PushBackNamed(NamedHandler{Name: "awssdk.client.LogResponse", Fn: logResponse})
}
