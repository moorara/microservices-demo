# sensor-service

## API

| Method   | Endpoint                 | Status | Response           | Description                      |
|----------|--------------------------|:------:|--------------------|----------------------------------|
| `POST`   | `/v1/sensors`            | `201`  | `sensor object`    | Creates a new sensor for a site  |
| `GET`    | `/v1/sensors?siteId=:id` | `200`  | `array of sensors` | Retrieves all sensors for a site |
| `GET`    | `/v1/sensors/:id`        | `200`  | `sensor object`    | Retrieves an existing sensor     |
| `PUT`    | `/v1/sensors/:id`        | `204`  |                    | Updates an existing sensor       |
| `DELETE` | `/v1/sensors/:id`        | `204`  |                    | Deletes an existing sensor       |

### Examples

```bash
curl \
  -H 'Content-Type: application/json' \
  -X POST \
  -d '{"siteId":"1111-aaaa","name":"temperature","unit":"celsius","minSafe":-30.0,"maxSafe":30.0}' \
  http://localhost:4020/v1/sensors

curl \
  -H 'Content-Type: application/json' \
  -X GET \
  http://localhost:4020/v1/sensors?siteId=1111-aaaa

curl \
  -H 'Content-Type: application/json' \
  -X GET \
  http://localhost:4020/v1/sensors/:id

curl \
  -H 'Content-Type: application/json' \
  -X PUT \
  -d '{"siteId":"1111-aaaa","name":"temperature","unit":"farenheit","minSafe":-22.0,"maxSafe":86.0}' \
  http://localhost:4020/v1/sensors/:id

curl \
  -H 'Content-Type: application/json' \
  -X DELETE \
  http://localhost:4020/v1/sensors/:id
```

## Commands

| Command                      | Description                             |
|------------------------------|-----------------------------------------|
| `make run`                   | Run the service locally                 |
| `make build`                 | Build the service binary locally        |
| `make test`                  | Run the unit tests                      |
| `make coverage`              | Run the unit tests with coverage report |
| `make docker`                | Build Docker image                      |
| `make docker-test`           | Build Docker image for testing          |
| `make push`                  | Push built images to registry           |
| `make save-docker`           | Save built images to disk               |
| `make load-docker`           | Load saved images from disk             |
| `make up`                    | Run the service locally in containers   |
| `make down`                  | Stop and removes local containers       |
| `make test-component`        | Run the component tests                 |
| `make test-component-docker` | Run the component tests in containers   |

## Documentation

  - https://gokit.io
  - http://www.gorillatoolkit.org/pkg/mux
  - https://godoc.org/github.com/prometheus/client_golang/prometheus
  - https://godoc.org/github.com/stretchr/testify

### OpenTracing & Jaeger

  - https://github.com/opentracing/opentracing-go
    - https://godoc.org/github.com/opentracing/opentracing-go
    - https://godoc.org/github.com/opentracing/opentracing-go/mocktracer
  - https://github.com/jaegertracing/jaeger-client-go
    - https://godoc.org/github.com/uber/jaeger-client-go
    - https://godoc.org/github.com/uber/jaeger-client-go/config
  - https://github.com/yurishkuro/opentracing-tutorial/tree/master/go
