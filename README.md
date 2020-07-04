# Implementación de SEArch

## Componentes

Por ahora habrá estos componentes en la arquitectura:

- broker
- middleware p/requires point
- middleware p/provides point
- repository


## Comunicación entre componentes

gRPC 

Me permite utilizar ProtoBuf para definir los tipos de los mensajes, su encoding y serialización en un stream de bits, definir las signaturas de los mensajes RPC (no define coreografías).

## Lenguaje de programación elegido para cada componente

Por qué elijo el lenguaje X?

Candidatos:

- go
- rust
- python

- erlang
- jolie

## Alcance de un MVP

No sería necesario monitoreo del estado de la interacción. 
