package imds

import (
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type APIHandlers interface {
	GetAPITokenHandler() http.Handler
	GetAPIHandler() http.Handler
}

func newTestServeMux(t *testing.T, handlers APIHandlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle(getTokenPath, validateAPITokenRequest(t, handlers.GetAPITokenHandler()))
	mux.Handle("/latest/", handlers.GetAPIHandler())

	return mux
}

func validateAPITokenRequest(t *testing.T, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if e, a := "PUT", r.Method; e != a {
			t.Errorf("expect %v, http method got %v", e, a)
			http.Error(w, http.StatusText(400), 400)
			return
		}
		if len(r.Header.Get(tokenTTLHeader)) == 0 {
			t.Errorf("expect token TTL header to be present in the request headers, got none")
			http.Error(w, http.StatusText(400), 400)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

type secureAPIHandler struct {
	t *testing.T

	tokens     []string
	tokenTTL   time.Duration
	apiHandler http.Handler

	activeToken atomic.Value
}

func newSecureAPIHandler(t *testing.T, tokens []string, tokenTTL time.Duration, apiHandler http.Handler) *secureAPIHandler {
	return &secureAPIHandler{
		t:          t,
		tokens:     tokens,
		tokenTTL:   tokenTTL,
		apiHandler: apiHandler,
	}
}

func (h *secureAPIHandler) GetAPITokenHandler() http.Handler {
	return http.HandlerFunc(h.handleAPIToken)
}

func (h *secureAPIHandler) handleAPIToken(w http.ResponseWriter, r *http.Request) {
	token := h.tokens[0]

	// set the active token
	h.storeActiveToken(token)

	// rotate the token
	if len(h.tokens) > 1 {
		h.tokens = h.tokens[1:]
	}

	var tokenTTLHeaderVal string
	if h.tokenTTL == 0 {
		tokenTTLHeaderVal = r.Header.Get(tokenTTLHeader)
	} else {
		tokenTTLHeaderVal = strconv.Itoa(int(h.tokenTTL / time.Second))
	}

	// set the header and response body
	w.Header().Set(tokenTTLHeader, tokenTTLHeaderVal)
	activeToken := h.getActiveToken()

	w.Write([]byte(activeToken))
}

func (h *secureAPIHandler) GetAPIHandler() http.Handler {
	return http.HandlerFunc(h.handleAPI)
}

func (h *secureAPIHandler) handleAPI(w http.ResponseWriter, r *http.Request) {
	token := h.getActiveToken()
	if len(token) == 0 {
		h.t.Errorf("expect token to have been requested, was not")
		http.Error(w, http.StatusText(401), 401)
		return
	}

	if e, a := token, r.Header.Get(tokenHeader); e != a {
		h.t.Errorf("expect %v token, got %v", e, a)
		http.Error(w, http.StatusText(401), 401)
		return
	}

	// delegate to configure handler for the request
	h.apiHandler.ServeHTTP(w, r)
}

func (h *secureAPIHandler) storeActiveToken(t string) {
	h.activeToken.Store(t)
}

func (h *secureAPIHandler) getActiveToken() string {
	activeToken := h.activeToken.Load()
	v, ok := activeToken.(string)
	if !ok {
		h.t.Errorf("expect valid active token string, got %T, %v", v, v)
	}

	return v
}

type insecureAPIHandler struct {
	t               *testing.T
	apiTokenErrCode int
	apiHandler      http.Handler
}

func newInsecureAPIHandler(t *testing.T, apiTokenErrCode int, apiHandler http.Handler) *insecureAPIHandler {
	return &insecureAPIHandler{
		t:               t,
		apiTokenErrCode: apiTokenErrCode,
		apiHandler:      apiHandler,
	}
}

func (h *insecureAPIHandler) GetAPITokenHandler() http.Handler {
	return http.HandlerFunc(h.handleAPIToken)
}

func (h *insecureAPIHandler) handleAPIToken(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(h.apiTokenErrCode), h.apiTokenErrCode)
}

func (h *insecureAPIHandler) GetAPIHandler() http.Handler {
	return http.HandlerFunc(h.handleAPI)
}

func (h *insecureAPIHandler) handleAPI(w http.ResponseWriter, r *http.Request) {
	if len(r.Header.Get(tokenHeader)) != 0 {
		h.t.Errorf("request token found, expected none")
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// delegate to configure handler for the request
	h.apiHandler.ServeHTTP(w, r)
}

type unauthorizedAPIHandler struct {
	t *testing.T

	enabled          bool
	secureAPIHandler *secureAPIHandler
}

func newUnauthorizedAPIHandler(t *testing.T, secureHandler *secureAPIHandler) *unauthorizedAPIHandler {
	return &unauthorizedAPIHandler{
		t:                t,
		secureAPIHandler: secureHandler,
	}
}

func (h *unauthorizedAPIHandler) GetAPITokenHandler() http.Handler {
	return http.HandlerFunc(h.handleAPIToken)
}

func (h *unauthorizedAPIHandler) handleAPIToken(w http.ResponseWriter, r *http.Request) {
	// Respond with 404 first, then token after 401 API handler response
	if !h.enabled {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	h.secureAPIHandler.GetAPITokenHandler().ServeHTTP(w, r)
}

func (h *unauthorizedAPIHandler) GetAPIHandler() http.Handler {
	return http.HandlerFunc(h.handleAPI)
}

func (h *unauthorizedAPIHandler) handleAPI(w http.ResponseWriter, r *http.Request) {
	// Respond with 401 first, then 200 for second. When enabled switch to
	// secure flow.
	if !h.enabled {
		h.enabled = true
		http.Error(w, http.StatusText(401), 401)
		return
	}

	h.secureAPIHandler.GetAPIHandler().ServeHTTP(w, r)
}

type requestTrace struct {
	requests []string
	mu       sync.Mutex
}

func newRequestTrace() *requestTrace {
	return &requestTrace{}
}

func (t *requestTrace) WrapHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.mu.Lock()
		t.requests = append(t.requests, r.URL.Path)
		t.mu.Unlock()

		handler.ServeHTTP(w, r)
	})
}
