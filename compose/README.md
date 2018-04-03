# docker-compose

## Commands

| Command            | Description                                          |
|--------------------|------------------------------------------------------|
| `make pull`        | Pulls required Docker images                         |
| `make images`      | Builds custom Docker images                          |
| `make services`    | Builds Docker images for application services        |
| `make up`          | Brings up a local environment using `docker-compose` |
| `make down`        | Takes down the local environment containers          |
| `make clean`       | Removes created Docker volumes                       |
| `test-up`          | Brings up a subset of local environment for testing  |
| `test-integration` | Runs the integration tests                           |
| `init-data`        | Initializes databases with sample data               |

## Dashboards

| Dashboard     | URL                                            | Required Information                 |
|---------------|------------------------------------------------|--------------------------------------|
| Kibana        | [http://localhost:5601](http://localhost:5601) | Index Pattern: `fluentd`             |
| Grafana       | [http://localhost:3000](http://localhost:3000) | User: `admin` <br/> Password: `pass` |
| Prometheus    | [http://localhost:9090](http://localhost:9090) | -                                    |
| Alert Manager | [http://localhost:9093](http://localhost:9093) | -                                    |
| cAdvisor      | [http://localhost:9080](http://localhost:9080) | -                                    |
| Træfik        | [http://localhost:2080](http://localhost:2080) | -                                    |

## Guides

### Resources

  * **Docker**
    - https://docs.docker.com/engine/reference/builder
    - https://docs.docker.com/compose/compose-file

  * **Fluentd**
    - https://docs.fluentd.org/v1.0/articles/life-of-a-fluentd-event
    - https://docs.fluentd.org/v1.0/articles/config-file
    - https://docs.fluentd.org/v1.0/articles/logging

  * **ElasticSearch**
    - https://www.elastic.co/guide/en/elasticsearch/reference/current/important-settings.html

  * **Kibana**
    - https://www.elastic.co/guide/en/kibana/current/settings.html

  * **Prometheus**
    - https://prometheus.io/docs/concepts
    - https://prometheus.io/docs/prometheus/latest/configuration/configuration
    - https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules
    - https://prometheus.io/docs/prometheus/latest/configuration/recording_rules

  * **Traefik**
    - https://docs.traefik.io/basics
    - https://docs.traefik.io/user-guide/examples
    - https://docs.traefik.io/configuration/commons
    - https://docs.traefik.io/configuration/api
    - https://docs.traefik.io/configuration/metrics
    - https://docs.traefik.io/configuration/entrypoints
    - https://docs.traefik.io/configuration/backends/rest
    - https://docs.traefik.io/configuration/backends/docker
