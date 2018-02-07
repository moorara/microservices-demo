# go-service

## API

| Verb     | Endpoint               | Description                    |
|----------|------------------------|--------------------------------|
| `POST`   | `/v1/votes`            | Creates a new vote for a link  |
| `GET`    | `/v1/votes?linkId=:id` | Retrieves all votes for a link |
| `GET`    | `/v1/votes/:id`        | Retrieves an existing vote     |
| `DELETE` | `/v1/votes/:id`        | Deletes an existing vote       |

### Examples

```bash
curl -H 'Content-Type: application/json' -X POST -d '{"linkId":"<linkId>", "stars":5}' http://localhost:4010/v1/votes
curl -H 'Content-Type: application/json' -X GET http://localhost:4010/v1/votes/<voteId>
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
