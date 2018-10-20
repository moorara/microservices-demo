# site-service

## API

| Method   | Endpoint        | Status | Response         | Description         |
|----------|-----------------|:------:|------------------|---------------------|
| `POST`   | `/v1/sites`     | `201`  | `site object`    | Creates a new site  |
| `GET`    | `/v1/sites`     | `200`  | `array of sites` | Retrieves all sites |
| `GET`    | `/v1/sites/:id` | `200`  | `site object`    | Retrieves a site    |
| `PUT`    | `/v1/sites/:id` | `204`  |                  | Updates a site      |
| `PATCH`  | `/v1/sites/:id` | `200`  | `site object`    | Modifies a site     |
| `DELETE` | `/v1/sites/:id` | `200`  | `site object`    | Removes a site      |

### Examples

```bash
curl \
  -H 'Content-Type: application/json' \
  -X POST \
  -d '{"name":"plant","location":"ottawa","tags":["power"],"priority":3}' \
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
  -d '{"name":"plant site","location":"toronto","tags":["power","hydro"],"priority":2}' \
  http://localhost:4010/v1/sites/:id

curl \
  -H 'Content-Type: application/json' \
  -X PATCH \
  -d '{"location":"kingston","priority":4}' \
  http://localhost:4010/v1/sites/:id

curl \
  -H 'Content-Type: application/json' \
  -X DELETE \
  http://localhost:4010/v1/sites/:id
```

## Commands

| Command                      | Description                             |
|------------------------------|-----------------------------------------|
| `yarn start`                 | Run the service locally                 |
| `yarn run nsp`               | Identify known vulneberities in service |
| `yarn run lint`              | Run standard linter                     |
| `yarn run lint:fix`          | Fix linting issues                      |
| `yarn run test`              | Run the unit tests                      |
| `yarn run test:component`    | Run the component tests                 |
| `make docker`                | Build Docker image                      |
| `make docker-test`           | Build Docker image for testing          |
| `make up`                    | Run the service locally in containers   |
| `make down`                  | Stop and removes local containers       |
| `make test-docker`           | Run the unit tests in containers        |
| `make test-component-docker` | Run the component tests in containers   |
| `make push`                  | Push built images to registry           |
| `make save-images`           | Save built images to disk               |

## Documentation

  - https://lodash.com
  - https://expressjs.com
  - http://mongoosejs.com
  - https://github.com/winstonjs/winston
  - https://github.com/bithavoc/express-winston
  - https://github.com/siimon/prom-client
  - https://shouldjs.github.io
  - https://mochajs.org
  - http://sinonjs.org
  - https://standardjs.com

### OpenTracing & Jaeger

  - https://github.com/opentracing/opentracing-javascript
    - https://opentracing-javascript.surge.sh
  - https://github.com/jaegertracing/jaeger-client-node
  - https://github.com/yurishkuro/opentracing-tutorial/tree/master/nodejs
