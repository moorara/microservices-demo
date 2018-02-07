# node-service

## API

| Verb     | Endpoint        | Description         |
|----------|-----------------|---------------------|
| `POST`   | `/v1/links`     | Creates a new link  |
| `GET`    | `/v1/links`     | Retrieves all links |
| `GET`    | `/v1/links/:id` | Retrieves a link    |
| `PUT`    | `/v1/links/:id` | Updates a link      |
| `DELETE` | `/v1/links/:id` | Removes a link      |

### Examples

```bash
curl -H 'Content-Type: application/json' -X POST -d '{"url":"https://docker.com", "title":"Docker", "tags":["container"], "rank":1}' http://localhost:4020/v1/links
curl -H 'Content-Type: application/json' -X GET http://localhost:4020/v1/links
curl -H 'Content-Type: application/json' -X GET http://localhost:4020/v1/links/<linkId>
curl -H 'Content-Type: application/json' -X PUT -d '{"url":"https://kubernetes.io", "title":"Kubernetes", "tags":["container"]}' http://localhost:4020/v1/links/<linkId>
curl -H 'Content-Type: application/json' -X DELETE http://localhost:4020/v1/links/<linkId>
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
