FROM golang:1.11.2 as builder
#Create package folder
RUN mkdir /go/src/go-chat-backend
#Copy the source files into the container
COPY / /go/src/go-chat-backend
#Set the go path to /go
ENV GOPATH=/go
#Set the current directory to the go package
WORKDIR /go/src/go-chat-backend
#Install dependencies
RUN go get -v ./...
#Build the application
RUN CGO_ENABLED=0 go install -a --installsuffix cgo --ldflags="-s" ./...
#Open new empty container
FROM alpine:latest
#Make sure the following commands are run as root
USER root
#Create a non-privileged user
RUN adduser -D -g appuser appuser
#Create the /var/app folder
RUN mkdir -p /var/app
RUN chown -R appuser:appuser /var/app
#Install package certs
RUN apk --no-cache add ca-certificates
#Copy over server
COPY --from=builder --chown=appuser:appuser /go/bin/server /var/app/server
#Make the server executable
RUN chmod a+x /var/app/server
#Set User to un-privileged user
USER appuser
#Set working directory
WORKDIR /var/app/
#Expose the server port
EXPOSE 9090
#Set entrypoint
ENTRYPOINT ["/var/app/server"]
