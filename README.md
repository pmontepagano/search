# SEArch


## Quickstart example

Do you want to just see an example of the infrastructure runnung?

1. **Tools:** Docker (v25+) and Docker Compose (v2.24+).
2. Clone the [repository](https://github.com/pmontepagano/search.git) and checkout the branch `origin/bisimulation-impl-paper`, or unzip [the file containing it](https://github.com/pmontepagano/search/archive/refs/heads/bisimulation-impl-paper.zip).
3. Once in the repository root, navigate in your terminal to the `examples/credit-card-payments` directory: `cd examples/credit-card-payments/`
4. Run the client: `docker-compose run client`

Before running the client, Docker will first build the infraestructure (i.e., the Broker and the Middleware), the services and the client. Then, it will run everything within containers; the Broker, a Middleware for each service required by the example, and their corresponding services behind them. Finally, it will run a Middleware for the client, and the client application behind it.


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

If you want to modify and/or contribute to this project, you will need the following tools installed:


1. [Python 3.12](https://python.org/) and these Python packages:
  - [python-betterproto](https://github.com/danielgtaylor/python-betterproto), used to compile our Protocol Buffer definitions to Python code.
  - [cfsm-bisimulation](https://github.com/diegosenarruzza/bisimulation/), used to calculate compatibility between contracts.
  - [z3-solver](https://github.com/Z3Prover/z3), which is a dependency of `cfsm-bisimulation`.

A simple way of installing all of these is to install Python 3.12 and run `pip install -r requirements.txt`

2. [Go 1.21](https://go.dev/) and these tools
  - [buf](https://buf.build/docs/installation), used to generate code from our Protocol Buffer definitions.
  - [mockery](https://vektra.github.io/mockery/), used to generate mocks for tests.

### Code generation from Protocol Buffer definitions

Our Protocol Buffer definitions live in the `proto` directory. The generated code lives in the `gen` directory. If you modify the Protocol Buffer definitions or add more output languages (see `buf.gen.yaml`), just run:

    buf generate proto

The generated code should be commited to the repository.

When modifying the Protocol Buffers, make sure you run the linter like this:

    buf lint proto

### Code generation for database management

We use [Ent](https://github.com/ent/ent) to model the database where we store contracts, providers and compatibility checks already computed.

The database definitions live in `ent/schema`. The rest of the files and directories in the `ent` directory are auto-generated from the schema definitions. If you change the schemas, run the folllowing command to regenerate code and commit any changes to the repository:

    go generate ./ent

#### Other useful Ent commands

##### Show schema in CLI

    go run -mod=mod entgo.io/ent/cmd/ent describe ./ent/schema

##### Show schema in [Atlas Cloud](https://gh.atlasgo.cloud/)

    go run -mod=mod ariga.io/entviz ./ent/schema

##### Generate Entity Relation diagram locally

    go run -mod=mod github.com/a8m/enter ./ent/schema

### Code generation for test mocks

 We use [mockery](https://vektra.github.io/mockery/) to generate test _mocks_ (only for the `contract` package for now). If you modify the `contact` package, regenerate the mocks by running the following command and commiting the changes to the repository. These generated mocks live in the `mocks` directory:

    mockery --dir contract --all --with-expecter

### Run tests

    go test ./...

And with [race detector](https://go.dev/doc/articles/race_detector):

    go test ./... -count=1 -race

### To get a report of code coverage

    go test ./... -coverprofile=coverage.txt -covermode atomic -coverpkg=./cfsm/...,./internal/...,./contract -timeout 30s
    go tool cover -html=coverage.txt

### To compile broker and middleware binaries:

    go build -o . ./...

## The example

```
cd examples/credit-card-payments
docker-compose build 