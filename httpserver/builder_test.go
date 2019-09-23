package httpserver_test

import "net/http"

type noopHTTPHandler struct{}

func (h noopHTTPHandler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {}
