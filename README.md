# Authentication microservice

The goal of the project is to explore the microservice nature in the form of an authentication service that allows the client to receive JWT tokens and check their validity.

Every part of the service is independent of each other. All parts are using `TokenService` interface. It allows to inject a new layer like logging, transport layer, different business logic, cache or something else easily.

## What is the project consists of
- HTTP and GRPC servers
- HTTP and GRPC clients
- TLS connection
- JWT and Bare Token Service

## Import clients

Import with `go get`

```
go get github.com/danblok/auth/client
```

## How to use
Make sure you have the `.env` file in your root directory with env vars like in `.env.example`.
Start with Docker Compose
```
docker compose up
```
or on your machine
```
make run
```
## Usefull data

`data` directory contains certificates and keys. It is possible to regenerate these keys
```
data/genkeys.sh
```
