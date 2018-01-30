# Microservices Toys

# go-service

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

# node-service

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

# compose

| Command         | Description                                          |
|-----------------|------------------------------------------------------|
| `make pull`     | Pulls required Docker images                         |
| `make fluentd`  | Builds custom `fluentd` Docker image                 |
| `make services` | Builds Docker images for application services        |
| `make up`       | Brings up a local environment using `docker-compose` |
| `make down`     | Takes down the local environment containers          |
| `make clean`    | Removes created Docker volumes                       |

The following web applications can be accessed:

| Web UI        | URL                                            | Required Information                 |
|---------------|------------------------------------------------|--------------------------------------|
| Kibana        | [http://localhost:5601](http://localhost:5601) | Index Pattern: `fluentd`             |
| Grafana       | [http://localhost:3000](http://localhost:3000) | User: `admin` <br/> Password: `pass` |
| Prometheus    | [http://localhost:9090](http://localhost:9090) | -                                    |
| Alert Manager | [http://localhost:9093](http://localhost:9093) | -                                    |
| cAdvisor      | [http://localhost:9080](http://localhost:9080) | -                                    |
| Træfik        | [http://localhost:2080](http://localhost:2080) | -                                    |
