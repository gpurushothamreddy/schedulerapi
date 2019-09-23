package httpserver

// NewHealthCheckEndpointBuilder to setting up a new health check endpoint
func NewHealthCheckEndpointBuilder() EndpointBuilder {
	return &getHealthCheckEndpointBuilder{}
}

type getHealthCheckEndpointBuilder struct {
}

// Build the endpoint
func (builder *getHealthCheckEndpointBuilder) Build() Endpoint {

	handler := NewHealthcheckHandler()

	return NewEndpoint(
		handler,
	)
}
