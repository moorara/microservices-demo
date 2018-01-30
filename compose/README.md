# docker-compose

## Commands

| Command         | Description                                          |
|-----------------|------------------------------------------------------|
| `make pull`     | Pulls required Docker images                         |
| `make fluentd`  | Builds custom `fluentd` Docker image                 |
| `make services` | Builds Docker images for application services        |
| `make up`       | Brings up a local environment using `docker-compose` |
| `make down`     | Takes down the local environment containers          |
| `make clean`    | Removes created Docker volumes                       |

## Dashboards

| Dashboard     | URL                                            | Required Information                 |
|---------------|------------------------------------------------|--------------------------------------|
| Kibana        | [http://localhost:5601](http://localhost:5601) | Index Pattern: `fluentd`             |
| Grafana       | [http://localhost:3000](http://localhost:3000) | User: `admin` <br/> Password: `pass` |
| Prometheus    | [http://localhost:9090](http://localhost:9090) | -                                    |
| Alert Manager | [http://localhost:9093](http://localhost:9093) | -                                    |
| cAdvisor      | [http://localhost:9080](http://localhost:9080) | -                                    |
| Træfik        | [http://localhost:2080](http://localhost:2080) | -                                    |
