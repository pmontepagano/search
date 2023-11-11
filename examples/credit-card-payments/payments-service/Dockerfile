ARG GOVERSION="1.21"
ARG USERNAME=search
FROM golang:${GOVERSION} as dev

RUN GRPC_HEALTH_PROBE_VERSION=v0.4.22 && \
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
COPY examples/credit-card-payments/payments-service ./

RUN go build -v -o /usr/local/bin/payments-service main.go

CMD ["/usr/local/bin/payments-service"]
