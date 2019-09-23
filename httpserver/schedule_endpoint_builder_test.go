package httpserver_test

import (
	"testing"

	. "github.com/onsi/gomega"
	httpserver "github.com/schedulerapi/httpserver"
	s "github.com/schedulerapi/scheduler"
)

func TestScheduleEndpointBuilder(t *testing.T) {
	g := NewGomegaWithT(t)
	t.Logf("%s - %s\n", "TestScheduleEndpointBuilder", "Building a http handler for the post schedule endpoint")
	var scheduler *s.Scheduler
	subject := httpserver.NewSchedulerEndpointBuilder(scheduler)

	endpoint := subject.Build()
	t.Logf("Should match expected http handler")
	g.Expect(endpoint.Handler()).To(Equal(httpserver.NewScheduleHandler(scheduler)))
}
