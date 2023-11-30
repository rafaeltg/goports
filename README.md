# goports
Service responsible for managing sea ports data.

It is comprised of two services:Contains two applications for managing ports (that have been defined via a JSON file):
* A service which exposes a REST API for managin ports (read, create and update)
* An ingestor which will parse a given JSON file, extract the port information, and call the service to persist the data

## Prerequisites
* `go` 1.21
* `docker`
* `make`

## Usage

### Starting applications
#### Server
The server can be started via Docker using:
```bash
docker-compose up -d server
```

The REST API will be exposed locally on port `8088`.

#### Ingestor
The ingestor can be started via Docker using:
```bash
docker-compose up -d ingestor
```

### Testing and analasying
The code can be linted and tested using the provided `Makefile`:
```bash
make lint
make test
```

## TODOs
- [] Add integration tests