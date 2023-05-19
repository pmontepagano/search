# Implementación de SEArch

## Componentes

Por ahora habrá estos componentes en la arquitectura:

- broker + repository (posiblemente después lo parta en dos)
- client middleware (está en un requires point)
- server middleware (está en un provides point)

## Código generado (go protobufs y go-grpc)

# Buf para protobufs y gRPC

  buf generate proto

# Entgo para manejo de la base de datos del broker

  go generate ./ent


## Comunicación entre componentes

### gRPC 

Me permite utilizar ProtoBuf para definir los tipos de los mensajes, su encoding y serialización en un stream de bits, definir las signaturas de los mensajes RPC (no define coreografías).


## Cómo ejecutar para entorno de desarrollo

Alcanza con tener Go instalado y ejecutar:

    go run broker/broker.go

En otra terminal:


    go run clientmiddleware/clientmiddleware.go


Y en otra:


    go run providermiddleware/providermiddleware.go


## Organización del código


En el directorio `protobuf` se encuentan los archivos `.proto` donde definimos los tipos de mensajes y los servicios. Esos archivos se compilan con [protoc](https://developers.google.com/protocol-buffers/docs/overview) y generan los archivos `.pb.go`.


# Building in Docker

- https://www.docker.com/blog/containerize-your-go-developer-environment-part-1/
- https://www.docker.com/blog/docker-golang/
