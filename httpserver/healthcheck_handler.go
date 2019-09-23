package httpserver

import (
	"encoding/json"
	"net/http"
)

type gethealthcheckHandler struct{}

// NewHealthcheckHandler for health check endpoint
func NewHealthcheckHandler() http.Handler {
	return &gethealthcheckHandler{}
}

//HealthcheckHandler
func (h *gethealthcheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(`{"status":"ok"}`)
}
