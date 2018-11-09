# Asset

**asset-service** is a microservice that creates, reads, updates, and deletes `assets`.

## TL;DR;

```bash
$ helm install repo/asset
```

## Introduction

This chart bootstraps a **asset-service** deployment on a Kubernetes cluster using the **Helm** package manager.

## Prerequisites

- Kubernetes 1.8+

## Installing the Chart

To install the chart with the release name `my-release`:

```bash
$ helm install --name my-release repo/asset
```

You can install or update the chart with the release name `my-release`:

```bash
$ helm upgrade --install my-release repo/asset
```

The command deploys **asset-service** on the Kubernetes cluster in the default configuration.
The configuration section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```bash
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the **asset** chart and their default values.

| Parameter                       | Description                                                 | Default                      |
|---------------------------------|-------------------------------------------------------------|------------------------------|
| `imagePullSecrets`              | A list of `docker-registry` secret names                    | `nil`                        |
| `image.repository`              | Docker image name                                           | `moorara/asset-service`      |
| `image.tag`                     | Docker image tag                                            | `0.1.0`                      |
| `image.pullPolicy`              | Docker image pull policy                                    | `IfNotPresent`               |
| `config.port`                   | Service http port                                           | `4040`                       |
| `config.logLevel`               | Service logging level                                       | `info`                       |
| `config.logSpans`               | Whether to log Jaeger spans                                 | `false`                      |
| `serviceAccount.create`         | Create Kubernetes service account for pod                   | `false`                      |
| `serviceAccount.name`           | The name of Kubernetes service account                      | `nil`                        |
| `pod.securityContext`           | The Kubernetes security context for pod                     | `{}`                         | 
| `pod.annotations`               | Kubernetes pod annotations                                  | `{}`                         |
| `deployment.replicaCount`       | Number of service replicas (pods)                           | `1`                          |
| `deployment.strategy`           | Deployment strategy for updating pods                       | `RollingUpdate`              |
| `deployment.annotations`        | Kubernetes deployment annotations                           | `{}`                         |
| `resources`                     | Resource requests and limits                                | `{}`                         |
| `nats.enabled`                  | Enable `nats` as a chart dependency                         | `true`                       |
| `nats.replicaCount`             | Number of NATS nodes (pods)                                 | `1`                          |
| `nats.auth.enabled`             | Enable authentication for NATS                              | `true`                       |
| `nats.auth.user`                | NATS username                                               | `asset-service`              |
| `nats.auth.password`            | NATS password                                               | `password!`                  |
| `nats.external.servers`         | External NATS servers                                       | `nil`                        |
| `nats.external.user`            | External NATS user                                          | `nil`                        |
| `nats.external.password`        | External NATS password                                      | `nil`                        |
| `cockroachdb.enabled`           | Enable `cockroachdb` as a chart dependency                  | `true`                       |
| `cockroachdb.Replicas`          | Number of CockroachDB nodes (pods)                          | `1`                          |
| `cockroachdb.database`          | CockroachDB database name                                   | `assets`                     |
| `cockroachdb.user`              | CockroachDB user                                            | `asset-service`              |
| `cockroachdb.password`          | CockroachDB password                                        | `password!`                  |
| `cockroachdb.external.addr`     | External CockroachDB address                                | `nil`                        |
| `cockroachdb.external.database` | ExternalCockroachDB database name                           | `nil`                        |
| `cockroachdb.external.user`     | External CockroachDB user                                   | `nil`                        |
| `cockroachdb.external.password` | External CockroachDB password                               | `nil`                        |
| `jaeger.enabled`                | Enable `jaeger-agent` sidecar for tracing                   | `false`                      |
| `jaeger.collector.address`      | Host and port for `jaeger-collector`                        | `nil`                        |
| `jaeger.agent.image`            | Docker image repository for `jaeger-agent`                  | `jaegertracing/jaeger-agent` |
| `jaeger.agent.tag`              | Docker image tag for `jaeger-agent`                         | `latest`                     |
| `jaeger.agent.pullPolicy`       | Docker image pull policy for `jaeger-agent`                 | `Always`                     |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example,

```bash
$ helm install --name my-release --set image.tag=latest,image.pullPolicy=Always repo/asset
```

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example,

```bash
$ helm install --name my-release -f values.yaml repo/asset
```

> **Tip**: You can use the default [values.yaml](values.yaml)
