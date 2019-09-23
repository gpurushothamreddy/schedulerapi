package httpserver

import (
	"github.com/schedulerapi/scheduler"
)

// NewGetReviewEndpointBuilder setting up a new review summary endpoint
func NewSchedulerEndpointBuilder(
	scheduler *scheduler.Scheduler,
) EndpointBuilder {
	return &getSchedulerEndpointBuilder{
		scheduler: scheduler,
	}
}

type getSchedulerEndpointBuilder struct {
	scheduler *scheduler.Scheduler
}

func (builder *getSchedulerEndpointBuilder) Build() Endpoint {

	handler := NewScheduleHandler(builder.scheduler)

	return NewEndpoint(
		handler,
	)
}
