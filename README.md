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

