package httpserver_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/schedulerapi/httpserver"
)

func TestRoutes(t *testing.T) {
	//setup
	getHealthCheckHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	postSchedulerHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	getHealthCheckEndpoint := httpserver.NewEndpoint(getHealthCheckHandler)
	postSchedulerEndpoint := httpserver.NewEndpoint(postSchedulerHandler)

	subject := httpserver.NewRouter(
		postSchedulerEndpoint,
		getHealthCheckEndpoint,
	)

	testCases := []struct {
		Name           string
		Desc           string
		UriRequestPath string
		Method         string
		Subject        http.Handler
	}{
		{
			Name:           "schedule router",
			Desc:           "routes to provided http handler for post schedule endpoint",
			UriRequestPath: "/schedule",
			Method:         "POST",
			Subject:        subject,
		}, {
			Name:           "health check router",
			Desc:           "routes to provided http handler for get health check endpoint",
			UriRequestPath: "/status/basic",
			Method:         "GET",
			Subject:        subject,
		},
	}
	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			t.Logf("%s - %s\n", c.Name, c.Desc)

			responseWriter := httptest.NewRecorder()

			request, err := http.NewRequest(c.Method, c.UriRequestPath, nil)
			if err != nil {
				t.Fatal(fmt.Sprintf("Expected `%v`, got: `%v`", nil, err))
			}
			c.Subject.ServeHTTP(responseWriter, request)

			if responseWriter.Code != http.StatusTeapot {
				t.Fatal(fmt.Sprintf("Expected `%v`, got: `%v`", http.StatusTeapot, responseWriter.Code))
			}
		})
	}
}
