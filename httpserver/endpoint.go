package httpserver

import (
	"net/http"
)

// Endpoint to configure endpoint
type Endpoint struct {
	handler http.Handler
}

func (e Endpoint) Handler() http.Handler { return e.handler }

func NewEndpoint(handler http.Handler) Endpoint {
	return Endpoint{
		handler: handler,
	}
}
