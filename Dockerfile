# syntax = docker/dockerfile:1-experimental

FROM golang:1.14-alpine AS build
WORKDIR /src
ENV CGO_ENABLED=1
COPY go.* .
RUN go mod download
COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o /out/broker broker/broker.go
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o /out/clientmiddleware clientmiddleware/clientmiddleware.go
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o /out/servermiddleware servermiddleware/servermiddleware.go