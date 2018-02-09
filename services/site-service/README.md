# site-service

## API

| Method   | Endpoint        | Description         |
|----------|-----------------|---------------------|
| `POST`   | `/v1/sites`     | Creates a new site  |
| `GET`    | `/v1/sites`     | Retrieves all sites |
| `GET`    | `/v1/sites/:id` | Retrieves a site    |
| `PUT`    | `/v1/sites/:id` | Updates a site      |
| `DELETE` | `/v1/sites/:id` | Removes a site      |

### Examples

```bash
curl \
  -H 'Content-Type: application/json' \
  -X POST \
  -d '{"name":"plant","location":"here","tags":["power"],"priority":3}' \
  http://localhost:4010/v1/sites

curl \
  -H 'Content-Type: application/json' \
  -X GET \
  http://localhost:4010/v1/sites

curl \
  -H 'Content-Type: application/json' \
  -X GET \
  http://localhost:4010/v1/sites/:id

curl \
  -H 'Content-Type: application/json' \
  -X PUT \
  -d '{"name":"plant site","location":"there","tags":["power","hydro"],"priority":2}' \
  http://localhost:4010/v1/sites/:id

curl \
  -H 'Content-Type: application/json' \
  -X DELETE \
  http://localhost:4010/v1/sites/:id
```

## Commands

| Command                   | Description                                |
|---------------------------|--------------------------------------------|
| `yarn start`              | Runs the service locally                   |
| `yarn run nsp`            | Identifies known vulneberities in service  |
| `yarn run lint`           | Runs standard linter                       |
| `yarn run test`           | Runs the unit tests                        |
| `yarn run test-component` | Runs the component tests                   |
| `make docker`             | Builds Docker image                        |
| `make up`                 | Runs the service locally in containers     |
| `make down`               | Stops and removes local containers         |
| `make docker-test`        | Builds Docker test image                   |
| `make test`               | Runs the unit tests in containers          |
| `make test-component`     | Runs the component tests in containers     |
