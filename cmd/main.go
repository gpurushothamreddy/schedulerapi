package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	//flags "github.com/jessevdk/go-flags"

	h "github.com/schedulerapi/httpserver"
	"github.com/schedulerapi/scheduler"
	"github.com/schedulerapi/scheduler/storage"
)

func configureScheduleEndpoint(scheduler *scheduler.Scheduler) {
	args := parse()
	mux := http.NewServeMux()
	mux.Handle("/", routingHTTPServer(scheduler))
	serverPort := strconv.Itoa(args.Server.Port)
	wsServer := &http.Server{
		Addr:    ":" + serverPort,
		Handler: mux,
	}
	fmt.Println("Scheduler API Service listening on port:" + serverPort)
	err := wsServer.ListenAndServe()
	if err != nil {
		fmt.Println("could not start the scheduler api server", err)
		os.Exit(1)
	}
}

func routingHTTPServer(scheduler *scheduler.Scheduler) http.Handler {
	getSchedulerEndpointBuilder := h.NewSchedulerEndpointBuilder(
		scheduler,
	)
	getHealthcheckEndpointBuilder := h.NewHealthCheckEndpointBuilder()
	getSchedulerEndpoint := getSchedulerEndpointBuilder.Build()
	getHealthcheckEndpoint := getHealthcheckEndpointBuilder.Build()
	return h.NewRouter(
		getSchedulerEndpoint,
		getHealthcheckEndpoint,
	)
}
func main() {

	storage := storage.NewMemoryStorage()

	s := scheduler.New(storage)

	err := s.Start()
	if err != nil {
		fmt.Println("could not start scheduler:", err)
		os.Exit(1)
	}
	fmt.Println("Scheduler up and running....")
	configureScheduleEndpoint(&s)
}
