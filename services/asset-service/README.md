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
