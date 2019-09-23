package httpserver_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/onsi/gomega"

	httpserver "github.com/schedulerapi/httpserver"
	"github.com/schedulerapi/scheduler"
	"github.com/schedulerapi/scheduler/storage"
)

func TestSchedulerHandler(t *testing.T) {
	g := NewGomegaWithT(t)
	testCases := []struct {
		Name                        string
		Desc                        string
		ExpectedResponse            string
		ExpectedStatusCode          int
		ExpectedContentType         string
		MuxProcessedHttpRequestURL  string
		MuxProcessedHttpRequestBody string
		Error                       error
	}{{
		Name:                       "TestSchedulerHandler - with invalid http request and bad response from handler",
		Desc:                       "Returns invalid body in response",
		ExpectedStatusCode:         http.StatusBadRequest,
		ExpectedContentType:        "application/json",
		MuxProcessedHttpRequestURL: "schedule",
		ExpectedResponse:           `{"errorCode":"INVALID_BODY","message":"Invalid body"}`,
	}, {
		Name:                       "TestSchedulerHandler - with valid http request and success response",
		Desc:                       "Returns success status code with empty response body",
		ExpectedStatusCode:         http.StatusCreated,
		ExpectedContentType:        "application/json",
		MuxProcessedHttpRequestURL: "schedule",
		MuxProcessedHttpRequestBody: `{
			"command": "ls -lah",
			"scheduleDateTime": "2099-09-22 15:27:30 PST"
			}`,
		ExpectedResponse: "",
	}, {
		Name:                       "TestSchedulerHandler - with invalid schedule date in http request body and bad response from handler",
		Desc:                       "Returns invalid body in response",
		ExpectedStatusCode:         http.StatusBadRequest,
		ExpectedContentType:        "application/json",
		MuxProcessedHttpRequestURL: "schedule",
		MuxProcessedHttpRequestBody: `{
			"command": "ls -lah",
			"scheduleDateTime": "1999-09-22 15:27:30 PST"
			}`,
		ExpectedResponse: `{"errorCode":"INVALID_BODY","message":"Invalid body"}`,
	}}
	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			t.Logf("%s - %s\n", c.Name, c.Desc)

			storage := storage.NewMemoryStorage()
			s := scheduler.New(storage)
			s.Start()
			schedulerHandler := httpserver.NewScheduleHandler(&s)
			clientRequest, err := http.NewRequest("POST", c.MuxProcessedHttpRequestURL, bytes.NewBuffer([]byte(c.MuxProcessedHttpRequestBody)))
			g.Expect(err).NotTo(HaveOccurred())
			responseWriter := httptest.NewRecorder()
			schedulerHandler.ServeHTTP(responseWriter, clientRequest)
			g.Expect(strings.TrimSpace(responseWriter.Body.String())).To(Equal(c.ExpectedResponse))
			g.Expect(responseWriter.Code).To(Equal(c.ExpectedStatusCode))
			g.Expect(responseWriter.Header().Get("Content-Type")).To(Equal(c.ExpectedContentType))
			s.Stop()
		})
	}
}
