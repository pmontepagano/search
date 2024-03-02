# SEArch


## Quickstart example

Do you want to just see an example of the infrastructure running?

1. **Tools:** Docker (v25+) and Docker Compose (v2.24+).
2. Clone the [repository](https://github.com/pmontepagano/search.git) and checkout the branch `origin/bisimulation-impl-paper`, or unzip [the file containing it](https://github.com/pmontepagano/search/archive/refs/heads/bisimulation-impl-paper.zip).
3. Once in the repository root, navigate in your terminal to the `examples/credit-card-payments` directory: 
```
cd examples/credit-card-payments/
```
4. Run the client:
```
docker-compose run client
```
Before running the client, Docker will first build the infraestructure (i.e., the Broker and the Middleware), the services and the client. Then, it will run everything within containers; the Broker, a Middleware for each service required by the example, and their corresponding services behind them. Finally, it will run a Middleware for the client, and the client application behind it.

## Setting up a development environment

### Prerequisites
Modify and/or contributing to this project requires the following tools to be installed:

1. **Tools:** [Python 3.12](https://python.org/) and the following Python packages:
  - [python-betterproto](https://github.com/danielgtaylor/python-betterproto), used to compile our Protocol Buffer definitions to Python code.
  - [cfsm-bisimulation](https://github.com/diegosenarruzza/bisimulation/), used to calculate compatibility between contracts.
  - [z3-solver](https://github.com/Z3Prover/z3), which is a dependency of `cfsm-bisimulation`.

A simple way of installing all of these is to install Python 3.12 and then run: 
```
    pip install -r requirements.txt
```

2. [Go 1.21](https://go.dev/) and the following tools:
  - [buf](https://buf.build/docs/installation), used to generate code from our Protocol Buffer definitions to Go code.
  - [mockery](https://vektra.github.io/mockery/), used to generate mocks for tests.
  - [ent](https://github.com/ent/ent), used to model the database where we store contracts, providers and compatibility checks already computed.


### Obtain the code
Clone the [repository](https://github.com/pmontepagano/search.git) and checkout the branch `origin/bisimulation-impl-paper`; alternatively unzip the file containing it.

The structure of the repository is:
- `internal`: contains the code of the middleware and the broker, along with their tests.
- `cfsm`: contains a _fork_ of the [`nickng/cfsm`](https://github.com/nickng/cfsm) library that implements Communicating Finite State Machines. The main changes compared to the original are: the addition of a parser to read FSA files, some data structure changes to force the serialization of CFSMs to be deterministic, and the addition of using CFSM names instead of numerical indices to make the generated FSA files more human-readable.
- `cmd`: contains the executable code for the middleware and broker.
- `contract`: contains the interfaces of the global contract and the local contract, along with concrete implementations using CFSMs.
- `gen`: contains the files generated by buf after compiling the .proto files. The middleware and broker use the Go code, but we also store the code generated in other languages.
- `ent`: contains the code for managing the broker database (ORM). The entity and relationship definitions are located in the schema subdirectory. The rest of the code is generated using the Ent framework.
- `mocks`: test mocks of the contracts generated with mockery, to facilitate the testing of the broker.
- `proto`: contains the Protocol Buffer files with the message type and gRPC service definitions. These files are compiled using buf. The configuration file that determines the target languages for compilation is `buf.gen.yaml` (in the root of the repository).


### Code generation from Protocol Buffer definitions
We use [buf](https://buf.build/docs/installation) to generate definition for message types and gRPC services. Our Protocol Buffer definitions live in the `proto` directory. The generated code lives in the `gen` directory. If you modify the Protocol Buffer definitions or add more output languages (see `buf.gen.yaml`), just run:
```
    buf generate proto
```
The generated code should be commited to the repository. After modifying the Protocol Buffers, make sure you run the linter like this:
```
    buf lint proto
```

### Code generation for database management
We use [ent](https://github.com/ent/ent) to model the database where we store contracts, providers and compatibility checks already computed.

The database definitions live in `ent/schema`. The rest of the files and directories in the `ent` directory are auto-generated from the schema definitions. If you change the schemas, run the folllowing command to regenerate code and commit any changes to the repository:
```
    go generate ./ent
```

#### Other useful Ent commands
- Show schema in CLI:
```
    go run -mod=mod entgo.io/ent/cmd/ent describe ./ent/schema
```
- Show schema in [Atlas Cloud](https://gh.atlasgo.cloud/):
```
    go run -mod=mod ariga.io/entviz ./ent/schema
```
- Generate Entity Relation diagram locally
```
    go run -mod=mod github.com/a8m/enter ./ent/schema
```

### Code generation for test mocks
In software development, test mocks are simulated objects used in testing. They replace real objects or modules, allowing developers to isolate and test specific code independently. This means tests run faster and are more reliable because they aren't influenced by external factors. Mocks also offer predictable behavior, enabling developers to define how they respond to interactions and pinpoint potential issues within the code under test. Essentially, test mocks facilitate focused, efficient, and reliable unit testing, contributing to stronger software development.

We use [mockery](https://vektra.github.io/mockery/) to generate test _mocks_ (only for the `contract` package for now). If you modify the `contract` package, regenerate the mocks by running the following command and commiting the changes to the repository. The generated mocks live in the `mocks` directory:
```
    mockery --dir contract --all --with-expecter
```

### Compiling the broker and the middleware:
Building the infrastructure (i.e., the Broker and the Middleware) account to executing the following command:
```
    go build -o . ./...
```
It will compile the code found in `internal`, `cfsm`, `contract`, `gen` and `ent`. This will generate the binaries `broker` and `middleware` and placed in `cmd`. Both programs have a `--help` flag:
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

### Running the tests
Running the tests requires running the following command:
```
    go test ./...
```
Running the tests with [race detector](https://go.dev/doc/articles/race_detector) requires the use to the following command:
```
    go test ./... -count=1 -race
```
In order to get a report of code coverage after running the tests, use the following commands:
```
    go test ./... -coverprofile=coverage.txt -covermode atomic -coverpkg=./cfsm/...,./internal/...,./contract -timeout 30s
    go tool cover -html=coverage.txt
```

### Compiling application


