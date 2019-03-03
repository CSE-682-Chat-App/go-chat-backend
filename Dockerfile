#
# Build Container
#
FROM golang:1.11.5-alpine3.7 as builder
ENV CGO_ENABLED=0

RUN apk update
RUN apk --no-cache add ca-certificates
RUN update-ca-certificates
RUN apk add git

RUN mkdir -p /go/src/github.com/CSE-682-Chat-App/go-chat-backend
COPY / /go/src/github.com/CSE-682-Chat-App/go-chat-backend/
ENV GOPATH=/go

WORKDIR /go/src/github.com/CSE-682-Chat-App/go-chat-backend/

RUN go get -v ./...
RUN go test ./...
RUN CGO_ENABLED=0 go install -a --installsuffix cgo --ldflags="-s" ./...

FROM alpine:3.7

RUN adduser -D -g app -u 1000 app
RUN mkdir /var/app
RUN chown app:app /var/app

USER app

COPY --from=builder --chown=app:app /go/bin/chat-backend /var/app/chat-backend

EXPOSE 8080

ENTRYPOINT ["sh", "-c", "/var/app/chat-backend"]
