package httpserver

import (
	"net/http"

	"github.com/bmizerany/pat"
)

// NewRouter build the router
func NewRouter(
	scheduleEndpoint Endpoint,
	healthCheckEndpoint Endpoint,
) http.Handler {
	patternServeMux := pat.New()

	patternServeMux.Post(schedulePath, scheduleEndpoint.Handler())
	patternServeMux.Get(healthCheckPath, healthCheckEndpoint.Handler())

	return patternServeMux
}

const schedulePath = "/schedule"
const healthCheckPath = "/status/basic"
