package httpserver

// EndpointBuilder interface for building endpoints
//go:generate counterfeiter . EndpointBuilder
type EndpointBuilder interface {
	Build() Endpoint
}
