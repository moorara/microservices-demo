# switch-service

## API

## Commands

| Command                      | Description                                          |
|------------------------------|------------------------------------------------------|
| `make dep`                   | Installs and updates dependencies                    |
| `make proto`                 | Generates gRPC code from protocol buffers definition |
| `make run`                   | Runs the service locally                             |
| `make build`                 | Builds the service binary locally                    |
| `make docker`                | Builds Docker image                                  |
| `make docker-test`           | Builds Docker test image                             |
| `make save-images`           | Saves built Docker images                            |
| `make up`                    | Runs the service locally in containers               |
| `make down`                  | Stops and removes local containers                   |
| `make test`                  | Runs the unit tests                                  |
| `make coverage`              | Runs the unit tests with coverage report             |
| `make test-component`        | Runs the component tests                             |
| `make test-component-docker` | Runs the component tests completely in containers    |

## Documentation

  - https://gokit.io

### Protocol Buffers & gRPC

  - https://developers.google.com/protocol-buffers
  - https://developers.google.com/protocol-buffers/docs/proto3
  - https://developers.google.com/protocol-buffers/docs/reference/go-generated
  - https://github.com/golang/protobuf
    - https://godoc.org/github.com/golang/protobuf/proto
    - https://godoc.org/github.com/golang/protobuf/protoc-gen-go
  - https://grpc.io
  - https://grpc.io/docs/guides
    - https://grpc.io/docs/guides/concepts
  - https://grpc.io/docs/tutorials
    - https://grpc.io/docs/tutorials/basic/go
  - https://grpc.io/docs/reference
    - https://grpc.io/docs/reference/go/generated-code
  - https://github.com/grpc/grpc-go
    - https://godoc.org/google.golang.org/grpc
      - https://godoc.org/google.golang.org/grpc#Server.ServeHTTP
    - https://godoc.org/google.golang.org/grpc/metadata

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
