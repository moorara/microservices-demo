# sensor-service

## API

| Method   | Endpoint                 | Description                      |
|----------|--------------------------|----------------------------------|
| `POST`   | `/v1/sensors`            | Creates a new sensor for a site  |
| `GET`    | `/v1/sensors?siteId=:id` | Retrieves all sensors for a site |
| `GET`    | `/v1/sensors/:id`        | Retrieves an existing sensor     |
| `DELETE` | `/v1/sensors/:id`        | Deletes an existing sensor       |

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
  http://localhost:4020/v1/sensors/:id
```

## Commands

| Command               | Description                               |
|-----------------------|-------------------------------------------|
| `make dep`            | Installs and updates dependencies         |
| `make run`            | Runs the service locally                  |
| `make build`          | Builds the service binary locally         |
| `make docker`         | Builds Docker image                       |
| `make up`             | Runs the service locally in containers    |
| `make down`           | Stops and removes local containers        |
| `make test`           | Runs the unit tests                       |
| `make coverage`       | Runs the unit tests with coverage report  |
| `make test-component` | Runs the component tests                  |
