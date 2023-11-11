FROM node:20-alpine

# Set the working directory to /app
WORKDIR /app

# WARNING! Paths in the Dockerfile are relative to the root of the repository.
# We do that to be able to copy the generated stubs from the gen directory.
# The context is set in the docker-compose.yml file.
COPY examples/credit-card-payments/payments-service/package*.json .
RUN npm install
COPY examples/credit-card-payments/payments-service/ .
RUN rm lib
COPY gen/ts ./lib
# RUN npm build

CMD ["npm", "start"]
