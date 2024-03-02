# SEArch Example: credit card payments

This directory contains a full example of how to use SEArch to develop software. The example includes one Service Client and two Service Providers. We run one instance of the middleware for each component, as well as one instance of the broker.

## Running the Example

Please follow the instructions below to run the SEArch example:

- Make sure you have [Docker Compose installed](https://docs.docker.com/compose/install/).
- Open a terminal, navigate to this directory and run:


    docker compose run client

You will be prompted by the Service Client for different questions. The output should look like this:

```
Please select at least one book to purchase:
1. The Great Gatsby
2. To Kill a Mockingbird
3. 1984
4. Pride and Prejudice
5. Animal Farm
6. The Catcher in the Rye
7. Lord of the Flies
8. Brave New World
9. The Hobbit
10. The Fellowship of the Ring
Enter the number(s) of the book(s) you want to purchase (comma-separated): 4,10,6,2
Enter your shipping address: Balcarce 50
Purchase Information:
Selected Books:
- Pride and Prejudice
- The Fellowship of the Ring
- The Catcher in the Rye
- To Kill a Mockingbird
Shipping Address: Balcarce 50
Total amount: 410.0
Enter the credit card number: 4111111111111
Enter the credit card expiration date: 11/25
Enter the credit card security code: 123
Payment nonce: 1234567890
Purchase successful!
```


To see all the logs, open another terminal, navigate to this same directory, and run:

    docker compose logs -f

That command will show you the logs all the containers:

- the broker
- the payments-service (a Service Provider)
- the backend (another Service Provider)
- the client (Service Client)
- 3 different instances of the middleware, one for each service
