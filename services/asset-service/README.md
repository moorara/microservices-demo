# asset-service

## API

## Commands

| Command                        | Description                                          |
|--------------------------------|------------------------------------------------------|
| `make dep`                     | Installs and updates dependencies                    |
| `make run`                     | Runs the service locally                             |
| `make build`                   | Builds the service binary locally                    |
| `make docker`                  | Builds Docker image                                  |
| `make docker-test`             | Builds Docker image for testing                      |
| `make save-images`             | Saves built Docker images                            |
| `make up`                      | Runs the service locally in containers               |
| `make down`                    | Stops and removes local containers                   |
| `make test`                    | Runs the unit tests                                  |
| `make coverage`                | Runs the unit tests with coverage report             |
| `make test-integration`        | Runs the integration tests                           |
| `make test-integration-docker` | Runs the integration tests completely in containers  |
| `make test-component`          | Runs the component tests                             |
| `make test-component-docker`   | Runs the component tests completely in containers    |

## Documentation

  - https://gokit.io

### ArangoDB

  - https://www.arangodb.com
  - https://www.arangodb.com/documentation
    - https://docs.arangodb.com/latest/Manual
      - https://docs.arangodb.com/latest/Manual/Scalability
      - https://docs.arangodb.com/latest/Manual/DataModeling
      - https://docs.arangodb.com/latest/Manual/Graphs
      - https://docs.arangodb.com/latest/Manual/Foxx
    - https://docs.arangodb.com/latest/AQL
      - https://docs.arangodb.com/latest/AQL/Fundamentals
      - https://docs.arangodb.com/latest/AQL/Operations
      - https://docs.arangodb.com/latest/AQL/DataQueries.html
      - https://docs.arangodb.com/latest/AQL/Functions
      - https://docs.arangodb.com/latest/AQL/Graphs
      - https://docs.arangodb.com/latest/AQL/Examples
    - https://docs.arangodb.com/latest/HTTP
  - https://www.arangodb.com/arangodb-training-center
  - https://github.com/arangodb/go-driver
    - https://godoc.org/github.com/arangodb/go-driver
    - https://godoc.org/github.com/arangodb/go-driver/http

### NATS

  - https://www.nats.io/documentation
    - https://www.nats.io/documentation/concepts/nats-messaging
      - https://www.nats.io/documentation/concepts/nats-pub-sub
      - https://www.nats.io/documentation/concepts/nats-req-rep
      - https://www.nats.io/documentation/concepts/nats-queueing
    - https://www.nats.io/documentation/internals/nats-protocol-demo
      - https://www.nats.io/documentation/internals/nats-protocol
      - https://www.nats.io/documentation/internals/nats-server-protocol
    - https://www.nats.io/documentation/server/gnatsd-intro
      - https://www.nats.io/documentation/server/gnatsd-slow-consumers
    - https://www.nats.io/documentation/streaming/nats-streaming-intro
      - https://www.nats.io/documentation/streaming/nats-streaming-quickstart
      - https://www.nats.io/documentation/streaming/nats-streaming-protocol
    - https://github.com/nats-io/go-nats/tree/master/examples
      - https://www.nats.io/documentation/tutorials/nats-pub-sub
      - https://www.nats.io/documentation/tutorials/nats-req-rep
      - https://www.nats.io/documentation/tutorials/nats-queueing
      - https://www.nats.io/documentation/tutorials/nats-client-dev
      - https://www.nats.io/documentation/tutorials/nats-benchmarking
  - https://github.com/pires/kubernetes-nats-cluster
  - https://github.com/canhnt/k8s-nats-streaming

### Metrics & Prometheus

  - https://github.com/prometheus/client_golang
    - https://godoc.org/github.com/prometheus/client_golang/prometheus
    - https://godoc.org/github.com/prometheus/client_golang/prometheus/promhttp

### OpenTracing & Jaeger

  - https://github.com/opentracing/opentracing-go
    - https://godoc.org/github.com/opentracing/opentracing-go
    - https://godoc.org/github.com/opentracing/opentracing-go/mocktracer
  - https://github.com/jaegertracing/jaeger-client-go
    - https://godoc.org/github.com/uber/jaeger-client-go
    - https://godoc.org/github.com/uber/jaeger-client-go/config
  - https://github.com/yurishkuro/opentracing-tutorial/tree/master/go
