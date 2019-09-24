# schedulerapi
schedule a command using REST api

## Dependencies

To run the schedule api project, you will need the following:

- [golang](https://blog.golang.org/go1.11) version 1.11 installed and available in your $PATH
- [dep](https://github.com/golang/dep) for go vendoring (package management)

## Setup / Configuration
This repo must be cloned into the path $GOPATH/src/github.com/schedulerapi. See below for more details.

1. Get the code:
```
git clone https://github.com/gpurushothamreddy/schedulerapi.git

```
2. Install `dep` (package manager) & this repo's dependencies:
```
brew install dep
dep ensure
```
* `dep` is a go dependency management tool. This may be retired in favor of go modules, as modules is the approach sanctioned by the core go team at Google.

Once that is done, you should be ready to roll. Run tests to verify installation:
```
./etc/bin/test
```

## Build & Run
To launch ScheduleApi microservice locally (and hit it via HTTP):
For running app locally execute following command:
```
./etc/bin/run
```
For running app locally with docker image execute following command:
```
./etc/bin/build-run-docker-image
```

Browse the url - `http://localhost:8888/status/basic` for testing health check of schedule api locally in above both the cases.

Schedule api exposes following endpoints:

```
Health check endpoint: To check the health status of the api 
curl -X GET \
  http://localhost:8888/status/basic \
  -H 'Accept: */*' \
  -H 'Host: localhost:8888'
  
Schedule endpoint: To schedule a bash command
curl -X POST \
  http://localhost:8888/schedule \
  -H 'Accept: */*' \
  -H 'Content-Length: 66' \
  -H 'Content-Type: application/json' \
  -H 'Host: localhost:8888' \
  -H 'cache-control: no-cache' \
  -d '{
"command": "ls",
"scheduleDateTime": "2020-09-23 16:09:50 PST"
}' 

Note: scheduleDateTime value should be always greater then current datetime.
```
## Limitations:
Current uses only in memory store of schedule task. However there is provision of to use sqlite storage to store schedule task in distributed env.

## Testing
To Run all unit and acceptance tests run locally use:
```
./etc/bin/test
```
This will also generate code coverage report
