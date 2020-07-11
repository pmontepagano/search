GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all: test build
build: 
	protoc \
	  --go_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:. \
	  --go-grpc_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:. \
	  --go_opt=paths=source_relative \
	  --go-grpc_opt=paths=source_relative \
	  protobuf/broker.proto
	$(GOBUILD) -v
test: 
	$(GOTEST) -v 
clean: 
	$(GOCLEAN)
run:
	$(GOCMD) run broker.go
