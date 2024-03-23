# SEArch logs description

Logs, as *SEArch* itself, are under development and, in some sense, they must undergo a process of major redesign. In their current state, logs only meet the basic requirement of being a debugging tool for the infrastructure. Such a redesign requires elaborating of a conceptual hierarchy of logable events such that they can expose relevant aspects of the state of each component (i.e., components of the infrastructure such as the middlewares and the broker, services and applications running within the infrastructure). 

One key limitation of the current logging infrastructure is that it is done by streaming the messages through the standard output of each of the processes involved in executing the application within *SEArch*; this posses difficulties for understanding executions as they relay on quite intricate communication schemes.

Regardless future plans, as for today, logs mainly allow basic monitoring of communication messages between components. This is split in four categories:

1. **middleware - broker:** encompases messages sent or received between a middleware and the broker, for example:

*A Middleware logs it's view of a brokerage request to the Broker:*
```
     middleware-client:10000 - Requesting brokerage of contract
```

2. **service/client - middleware:** encompases messages exchanged between agents participating in an execution and their corresponding middlewares. As we mentioned before, not having a unified logging infrastructure exposses limitations, mainly related to the heterogeneity derived from the idiosyncracies of the different programming languages involved and their associated logging libraries. Examples of these events are: 

*A Middleware logs it's view of a message sent by a process running behind it:*
```
     Received CardDetailsWithTotalAmount: {232 22 22 110}
```

3. **middleware - middleware:** encompases the logging events associated to messages exchanged between middlewares, for example:

*A Middleware logs when it dispatches a message to a another middleware hosting a remote participant:*
```
     Sent message to remote Service Provider for channel cf4ea1f5-44d3-4b27-b489-64c424911611, participant PPS
```

4. **critical internal actions of the infrastructure's component:** these messages are used for monitoring the behaviour of the components by following their execution by tracking chosen internal actions, for example: 

*A Middleware logs when it receives a first outbound message from a participant:*
```
     Received first outbound message to send on channel cf4ea1f5-44d3-4b27-b489-64c424911611 for participant Srv. Opening connection to remote Service Provider.
```
*A Middleware starts the routine that periodically sends outgoing messages on an outbox buffer:*
```
     Started sender routine for channel cf4ea1f5-44d3-4b27-b489-64c424911611, participant Srv
```

Naturally, as *SEArch* is still under development, many debug logging entries can be found in the code. These entries are used to get a close look at critical parts of the code implementing the different components, both of the infrastructure, and the example provided with the implementation.

## Logs when running within Docker containers

As we mentioned before, *SEArch* has a basic logging infrastructure based on using the standard output stream of the processes. 

When running within a Docker container it is possible to obtain and visualize the logs by using built in capability of the `docker compose` pluggin.

For a basic usage of this infrastructure a user can try the following commands:

1. Visualizing the a time based composition of all the containers:
```bash
     docker compose logs -f
```
2. Visualizing the log of a specific container:
```bash
     docker compose logs [container_identifier]
```
