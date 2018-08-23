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

## API Gateways

| Type    | Transport | Gateway | URL                            |
| --------|-----------|---------|--------------------------------|
| REST    | HTTP      | Træfik  | http://localhost:1080/api/v1/  |
| REST    | HTTPS     | Træfik  | https://localhost:1443/api/v1/ |
| REST    | HTTPS     | Caddy   | https://localhost/api/v1/      |
| GraphQL | HTTPS     | Caddy   | https://localhost/graphql      |


## Dashboards

| Dashboard         | URL                    | Username | Password | Required Information             |
|-------------------|------------------------|----------|----------|----------------------------------|
| **GraphiQL**      | http://localhost:5000  |          |          |                                  |
| **Kibana**        | http://localhost:5601  |          |          | Index Pattern: `fluentd`         |
| **Grafana**       | http://localhost:3000  | `admin`  | `pass`   | Source: `http://prometheus:9090` |
| **Prometheus**    | http://localhost:9090  |          |          |                                  |
| **Alert Manager** | http://localhost:9093  |          |          |                                  |
| **cAdvisor**      | http://localhost:9800  |          |          |                                  |
| **Jaeger UI**     | http://localhost:16686 |          |          |                                  |
| **Træfik**        | http://localhost:1900  |          |          |                                  |

## Ports

| Port       | Container        | Description                                                      |
|------------|------------------|------------------------------------------------------------------|
| `80`       | `caddy`          | Caddy http port                                                  |
| `443`      | `caddy`          | Caddy https port                                                 |
| `9900`     | `caddy`          | Caddy Prometheus metrics                                         |
| `1080`     | `traefik`        | Træfik http port                                                 |
| `1443`     | `traefik`        | Træfik https port                                                |
| `1900`     | `traefik`        | Træfik dashboard                                                 |
| `4000`     | `react-client`   | react-client service                                             |
| `4010`     | `site-service`   | site-service service                                             |
| `4020`     | `sensor-service` | sensor-service service                                           |
| `4030`     | `switch-service` | switch-service service (grpc)                                    |
| `4031`     | `switch-service` | switch-service service (http)                                    |
| `5000`     | `graphql-server` | graphql-server service                                           |
| `6379`     | `redis`          | Redis service                                                    |
| `27017`    | `mongo`          | MongoDB service                                                  |
| `8529`     | `arango`         | ArangoDB service                                                 |
| `5432`     | `postgres`       | PostgreSQL service                                               |
| `9200`     | `elasticsearch`  | Elasticsearch RESTful API                                        |
| `9300`     | `elasticsearch`  | Elasticsearch transport protocol                                 |
| `5601`     | `kibana`         | Kibana dashboard                                                 |
| `24224`    | `fluentd`        | Fluentd tcp and udp protocol                                     |
| `9090`     | `prometheus`     | Prometheus                                                       |
| `9091`     | `prometheus`     | Prometheus push gateway                                          |
| `3000`     | `grafana`        | Grafana dashboard                                                |
| `9093`     | `alertmanager`   | Alertmanager                                                     |
| `9100`     | `node-exporter`  | Prometheus node exporter                                         |
| `9800`     | `cadvisor`       | cAdvisor dashboard                                               |
| `5778`     | `jaeger`         | jaeger-agent: serve configurations, and sampling strategies      |
| `5775/udp` | `jaeger`         | jaeger-agent: accept zipkin.thrift over compact thrift protocol  |
| `6831/udp` | `jaeger`         | jaeger-agent: accept jaeger.thrift over compact thrift protocol  |
| `6832/udp` | `jaeger`         | jaeger-agent: accept jaeger.thrift over binary thrift protocol   |
| `9411`     | `jaeger`         | jaeger-collector: zipkin compatible endpoint                     |
| `14268`    | `jaeger`         | jaeger-collector: accept jaeger.thrift directly from clients     |
| `16686`    | `jaeger`         | jaeger-query: serve jaeger ui at `/` and api endpoints at `/api` |

## Documentation

### Docker

  - https://docs.docker.com/engine/reference/builder
  - https://docs.docker.com/compose/compose-file

### Fluentd

  - https://docs.fluentd.org/v1.0/articles/life-of-a-fluentd-event
  - https://docs.fluentd.org/v1.0/articles/config-file
  - https://docs.fluentd.org/v1.0/articles/logging

### ElasticSearch

  - https://www.elastic.co/guide/en/elasticsearch/reference/current/important-settings.html

### Kibana

  - https://www.elastic.co/guide/en/kibana/current/settings.html

### Prometheus

  - https://prometheus.io/docs/concepts
  - https://prometheus.io/docs/prometheus/latest/configuration/configuration
  - https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules
  - https://prometheus.io/docs/prometheus/latest/configuration/recording_rules

### Traefik

  - https://docs.traefik.io/basics
  - https://docs.traefik.io/user-guide/examples
  - https://docs.traefik.io/configuration/commons
  - https://docs.traefik.io/configuration/api
  - https://docs.traefik.io/configuration/metrics
  - https://docs.traefik.io/configuration/entrypoints
  - https://docs.traefik.io/configuration/backends/rest
  - https://docs.traefik.io/configuration/backends/docker

### Caddy

  - https://caddyserver.com/docs/http-caddyfile
  - https://caddyserver.com/docs/tls
  - https://caddyserver.com/docs/redir
  - https://caddyserver.com/docs/rewrite
  - https://caddyserver.com/docs/proxy
  
### OpenTracing & Jaeger

  - https://github.com/opentracing/specification
  - https://www.jaegertracing.io/docs/architecture
  - https://www.jaegertracing.io/docs/deployment
  - https://www.jaegertracing.io/docs/monitoring
