ARG GOVERSION="1.22"
ARG USERNAME=search
FROM golang:${GOVERSION} as dev

RUN GRPC_HEALTH_PROBE_VERSION=v0.4.24 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

WORKDIR /src
# ENV CGO_ENABLED=0
COPY go.* .
RUN go mod download
COPY cfsm ./cfsm
COPY contract ./contract
COPY ent ./ent
COPY mocks ./mocks
COPY gen/go ./gen/go
COPY cmd ./cmd
COPY internal ./internal

RUN --mount=type=cache,target=/home/$USERNAME/.cache/go-build \
    go build -v -o /usr/local/bin/broker cmd/broker/broker.go
RUN --mount=type=cache,target=/home/$USERNAME/.cache/go-build \
    go build -v -o /usr/local/bin/middleware cmd/middleware/middleware.go

# FROM scratch AS prod
# ENV PATH=/
# COPY --from=dev /bin/grpc_health_probe /
# COPY --from=dev /usr/local/bin/middleware /
# COPY --from=dev /usr/local/bin/broker /


# Use Python 3.12 as base image
FROM python:3.12 as with-python-vfsm-bisimulation
ENV PYTHONUNBUFFERED=1

# Set working directory to /app
WORKDIR /app

RUN pip install z3-solver cfsm-bisimulation

COPY --from=dev /bin/grpc_health_probe /usr/local/bin/
COPY --from=dev /usr/local/bin/middleware /usr/local/bin/
COPY --from=dev /usr/local/bin/broker /usr/local/bin/