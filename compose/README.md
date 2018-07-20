# docker-compose

## Commands

| Command            | Description                                          |
|--------------------|------------------------------------------------------|
| `make images`      | Builds custom Docker images                          |
| `make services`    | Builds Docker images for application services        |
| `make up`          | Brings up a local environment using `docker-compose` |
| `make down`        | Takes down the local environment containers          |
| `make clean`       | Removes created Docker volumes                       |
| `test-up`          | Brings up a subset of local environment for testing  |
| `test-integration` | Runs the integration tests                           |
| `init-data`        | Initializes databases with sample data               |

## Ports

| Port     | Container        | Description                      |
|----------|------------------|----------------------------------|
| `80`     | `caddy`          | Caddy http port                  |
| `443`    | `caddy`          | Caddy https port                 |
| `9900`   | `caddy`          | Caddy Prometheus metrics         |
| `1080`   | `traefik`        | Træfik http port                 |
| `1443`   | `traefik`        | Træfik https port                |
| `1900`   | `traefik`        | Træfik dashboard                 |
| `4000`   | `react-client`   | react-client service             |
| `4010`   | `site-service`   | site-service service             |
| `4020`   | `sensor-service` | sensor-service service           |
| `6379`   | `redis`          | Redis service                    |
| `27017`  | `mongo`          | MongoDB service                  |
| `8529`   | `arango`         | ArangoDB service                 |
| `5432`   | `postgres`       | PostgreSQL service               |
| `9200`   | `elasticsearch`  | Elasticsearch RESTful API        |
| `9300`   | `elasticsearch`  | Elasticsearch transport protocol |
| `5601`   | `kibana`         | Kibana dashboard                 |
| `24224`  | `fluentd`        | Fluentd tcp and udp protocol     |
| `9090`   | `prometheus`     | Prometheus                       |
| `9091`   | `prometheus`     | Prometheus push gateway          |
| `3000`   | `grafana`        | Grafana dashboard                |
| `9093`   | `alertmanager`   | Alertmanager                     |
| `9100`   | `node-exporter`  | Prometheus node exporter         |
| `9800`   | `cadvisor`       | cAdvisor dashboard               |

## Dashboards

| Dashboard                               | Required Information                                                        |
|-----------------------------------------|-----------------------------------------------------------------------------|
| [Kibana](http://localhost:5601)         | Index Pattern: `fluentd`                                                    |
| [Grafana](http://localhost:3000)        | User: `admin` <br/> Password: `pass` <br/> Source: `http://prometheus:9090` |
| [Prometheus](http://localhost:9090)     |                                                                             |
| [Alert Manager ](http://localhost:9093) |                                                                             |
| [cAdvisor](http://localhost:9800)       |                                                                             |
| [Træfik](http://localhost:1900)         |                                                                             |

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

  * **Caddy**
    - https://caddyserver.com/docs/http-caddyfile
    - https://caddyserver.com/docs/tls
    - https://caddyserver.com/docs/redir
    - https://caddyserver.com/docs/rewrite
    - https://caddyserver.com/docs/proxy
  
  * **OpenTracing & Jaeger**
    - https://github.com/opentracing/specification
    - https://www.jaegertracing.io/docs/architecture
    - https://www.jaegertracing.io/docs/deployment
    - https://www.jaegertracing.io/docs/monitoring
