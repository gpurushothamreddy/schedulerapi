FROM golang:alpine
RUN apk update && apk add bash 

COPY scheduleapi /scheduleapi

EXPOSE 8888

ENTRYPOINT ["/scheduleapi"]