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

| Command                      | Description                                       |
|------------------------------|---------------------------------------------------|
| `make dep`                   | Installs and updates dependencies                 |
| `make run`                   | Runs the service locally                          |
| `make build`                 | Builds the service binary locally                 |
| `make docker`                | Builds Docker image                               |
| `make docker-test`           | Builds Docker test image                          |
| `make save-images`           | Saves built Docker images                         |
| `make up`                    | Runs the service locally in containers            |
| `make down`                  | Stops and removes local containers                |
| `make test`                  | Runs the unit tests                               |
| `make coverage`              | Runs the unit tests with coverage report          |
| `make test-component`        | Runs the component tests                          |
| `make test-component-docker` | Runs the component tests completely in containers |
