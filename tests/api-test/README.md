# api-test

## Commands

| Command                        | Description                             |
|--------------------------------|-----------------------------------------|
| `make dep`                     | Install and updates dependencies        |
| `make run`                     | Run the service locally                 |
| `make build`                   | Build the service binary locally        |
| `make docker`                  | Build Docker image                      |
| `make test`                    | Run the unit tests                      |
| `make coverage`                | Run the unit tests with coverage report |
| `make push`                    | Push built images to registry           |
| `make save-images`             | Save built images to disk               |

## Documentation

  - https://gokit.io

### Metrics & Prometheus

  - https://github.com/prometheus/pushgateway
  - https://github.com/prometheus/client_golang
    - https://godoc.org/github.com/prometheus/client_golang/prometheus
    - https://godoc.org/github.com/prometheus/client_golang/prometheus/push
    - https://godoc.org/github.com/prometheus/client_golang/prometheus/promhttp

### OpenTracing & Jaeger

  - https://github.com/opentracing/opentracing-go
    - https://godoc.org/github.com/opentracing/opentracing-go
    - https://godoc.org/github.com/opentracing/opentracing-go/mocktracer
  - https://github.com/jaegertracing/jaeger-client-go
    - https://godoc.org/github.com/uber/jaeger-client-go
    - https://godoc.org/github.com/uber/jaeger-client-go/config
  - https://github.com/yurishkuro/opentracing-tutorial/tree/master/go
