# asset-service

## API

## Commands

| Command                        | Description                             |
|--------------------------------|-----------------------------------------|
| `make run`                     | Run the service locally                 |
| `make build`                   | Build the service binary locally        |
| `make docker`                  | Build Docker image                      |
| `make docker-test`             | Build Docker image for testing          |
| `make up`                      | Run the service locally in containers   |
| `make down`                    | Stop and removes local containers       |
| `make test`                    | Run the unit tests                      |
| `make coverage`                | Run the unit tests with coverage report |
| `make test-integration`        | Run the integration tests               |
| `make test-integration-docker` | Run the integration tests in containers |
| `make test-component`          | Run the component tests                 |
| `make test-component-docker`   | Run the component tests in containers   |
| `make push`                    | Push built images to registry           |
| `make save-images`             | Save built images to disk               |

## Documentation

  - https://gokit.io

### CockroachDB

  - https://www.cockroachlabs.com/docs/stable
    - https://www.cockroachlabs.com/docs/stable/install-cockroachdb.html
    - https://www.cockroachlabs.com/docs/stable/cockroach-commands.html
      - https://www.cockroachlabs.com/docs/stable/start-a-node.html
      - https://www.cockroachlabs.com/docs/stable/create-and-manage-users.html
    - https://www.cockroachlabs.com/docs/stable/start-a-local-cluster.html
    - https://www.cockroachlabs.com/docs/stable/secure-a-cluster.html
    - https://www.cockroachlabs.com/docs/stable/build-a-go-app-with-cockroachdb.html
    - https://www.cockroachlabs.com/docs/stable/build-a-go-app-with-cockroachdb-gorm.html

### GORM

  - http://gorm.io/docs
  - https://github.com/jinzhu/gorm
    - https://godoc.org/github.com/jinzhu/gorm

### NATS

  - https://www.nats.io/documentation
    - https://www.nats.io/documentation/concepts/nats-messaging
      - https://www.nats.io/documentation/concepts/nats-pub-sub
      - https://www.nats.io/documentation/concepts/nats-req-rep
      - https://www.nats.io/documentation/concepts/nats-queueing
    - https://www.nats.io/documentation/server/gnatsd-intro
      - https://www.nats.io/documentation/server/gnatsd-slow-consumers
    - https://www.nats.io/documentation/streaming/nats-streaming-intro
      - https://www.nats.io/documentation/streaming/nats-streaming-quickstart
    - https://github.com/nats-io/go-nats/tree/master/examples
      - https://www.nats.io/documentation/tutorials/nats-pub-sub
      - https://www.nats.io/documentation/tutorials/nats-req-rep
      - https://www.nats.io/documentation/tutorials/nats-queueing
      - https://www.nats.io/documentation/tutorials/nats-client-dev
      - https://www.nats.io/documentation/tutorials/nats-benchmarking
  - https://github.com/nats-io/nats-operator
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
