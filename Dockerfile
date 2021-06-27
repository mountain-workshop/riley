FROM golang:1.16 AS build-env
WORKDIR /go/src/github.com/mountain-workshop/riley

ARG version=dev
ARG commit=none

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN VERSION=$version COMMIT=$commit make build-linux

FROM alpine:3.13

USER nobody
COPY --from=build-env /go/src/github.com/mountain-workshop/riley /riley

ENTRYPOINT ["/riley"]
