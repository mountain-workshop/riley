FROM golang:1.16 AS build-env
WORKDIR /go/src/git.iratepublik.com/discord-house-cup/

ARG version=dev
ARG commit=none

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN VERSION=$version COMMIT=$commit make build-linux

FROM alpine:3.13

USER nobody
COPY --from=build-env /go/src/git.iratepublik.com/discord-house-cup/discord-house-cup /

ENTRYPOINT ["/discord-house-cup"]
