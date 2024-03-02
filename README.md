# SEArch


## Quickstart example

Do you want to just see an example of the infrastructure runnung?

You'll need Docker (v25+) and Docker Compose (v2.24+).

Navigate in your terminal to the `examples/credit-card-payments` directory and run:

```
cd examples/credit-card-payments/
docker-compose run client
```


## How to compile and run the Broker and Middleware

You need to have Go 1.21 installed and run:

    go build -o . ./...


This will generate the binaries `broker` and `middleware`. Both programs have a `--help` flag:

```
./broker --help
Usage of ./broker:
  -cert_file string
    	The TLS cert file
  -database_file string
    	The path for the broker database (default "SEARCHbroker.db")
  -host string
    	The host on which the broker listens (defaults to all interfaces)
  -key_file string
    	The TLS key file
  -port int
    	Port number on which the broker listens (default 10000)
  -tls
    	Connection uses TLS if true, else plain TCP
  -use_python_bisimulation
    	Use the Python Bisimulation library to check for bisimulation
```

```
./middleware --help
Usage of ./middleware:
  -broker_addr string
    	The server address in the format of host:port (default "localhost:")
  -cert_file string
    	The TLS cert file
  -key_file string
    	The TLS key file
  -private_host string
    	Host IP on which private service listens (default "localhost")
  -private_port int
    	The port for private services (default 11000)
  -public_host string
    	Host IP on which public service listens (defaults to all)
  -public_port int
    	The port for public facing middleware (default 10000)
  -public_url string
    	The URL for public facing middleware
  -tls
    	Connection uses TLS if true, else plain TCP
```

## How to see 

## How to set up a development environment

### Prerequisites

If you want to modify and/or contribute to this project, you will n

# Implementación de SEArch

## Código generado (go protobufs y go-grpc)

### Buf para protobufs y gRPC

- Tenemos que tener instalado [este paquete para generar los stubs de Python](https://github.com/danielgtaylor/python-betterproto). La manera más sencilla es crear un virtual environment de Python e instalar el archivo `requirements.txt` que se encuentra en la raíz de este repositorio.
- También tenemos que tener instalada la herramienta [buf](https://buf.build/docs/installation).
- Luego, basta con ejecutar el siguiente comando:

    buf generate proto

### Entgo para manejo de la base de datos del broker

    go generate ./ent

### Para regenerar mocks con [mockery](https://vektra.github.io/mockery/)

    mockery --dir contract --all --with-expecter

## Comandos varios

### To get a report of code coverage

    go test ./... -coverprofile=coverage.txt -covermode atomic -coverpkg=./cfsm/...,./internal/...,./contract -timeout 30s
    go tool cover -html=coverage.txt

### Para correr los tests

    go test ./...

Y con el [race detector](https://go.dev/doc/articles/race_detector):

    go test ./... -count=1 -race

### Para compilar los binarios de broker y middleware

    go build -o . ./...

### Comandos útiles de Entgo (ORM)

#### Show schema in CLI

    go run -mod=mod entgo.io/ent/cmd/ent describe ./ent/schema

#### Show schema in [Atlas Cloud](https://gh.atlasgo.cloud/)

    go run -mod=mod ariga.io/entviz ./ent/schema

#### Generate Entity Relation diagram locally

    go run -mod=mod github.com/a8m/enter ./ent/schema


## Run ChorGram's gc2fsa

After `--` you send the parameters. In this example, we simply pass the input file name.

    wasmtime --dir=. gc2fsa.wasm -- pingpong.gc

