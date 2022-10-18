# syntax = docker/dockerfile:1-experimental
# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.137.0/containers/go/.devcontainer/base.Dockerfile
ARG VARIANT="1.19"
ARG USERNAME=search
# You should use here your UID and GID
ARG USER_UID=501
ARG USER_GID=20
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT} as dev

WORKDIR /src
ENV CGO_ENABLED=0
COPY go.* .
RUN go mod download
COPY . .

RUN --mount=type=cache,target=/home/$USERNAME/.cache/go-build \
    go build -o /out/broker cmd/broker/broker.go
RUN --mount=type=cache,target=/home/$USERNAME/.cache/go-build \
    go build -o /out/middleware cmd/middleware/middleware.go

FROM scratch AS prod
COPY --from=dev /out/* /