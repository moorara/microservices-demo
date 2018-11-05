# Site

**site-service** is a microservice that creates, reads, updates, and deletes `sites`.

## TL;DR;

```bash
$ helm install repo/site
```

## Introduction

This chart bootstraps a **site-service** deployment on a Kubernetes cluster using the **Helm** package manager.

## Prerequisites

- Kubernetes 1.8+

## Installing the Chart

To install the chart with the release name `my-release`:

```bash
$ helm install --name my-release repo/site
```

You can install or update the chart with the release name `my-release`:

```bash
$ helm upgrade --install my-release repo/site
```

The command deploys **site-service** on the Kubernetes cluster in the default configuration.
The configuration section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```bash
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the **site** chart and their default values.

| Parameter                   | Description                                                        | Default                      |
|-----------------------------|--------------------------------------------------------------------|------------------------------|
| `imagePullSecrets`          | A list of `docker-registry` secret names                           | `nil`                        |
| `image.repository`          | Docker image name                                                  | `moorara/site-service`       |
| `image.tag`                 | Docker image tag                                                   | `0.1.0`                      |
| `image.pullPolicy`          | Docker image pull policy                                           | `IfNotPresent`               |
| `config.port`               | Service http port                                                  | `4010`                       |
| `config.logLevel`           | Service logging level                                              | `info`                       |
| `config.logSpans`           | Whether to log Jaeger spans                                        | `false`                      |
| `serviceAccount.create`     | Create Kubernetes service account for pod                          | `false`                      |
| `serviceAccount.name`       | The name of Kubernetes service account                             | `nil`                        |
| `pod.securityContext`       | The Kubernetes security context for pod                            | `{}`                         | 
| `pod.annotations`           | Kubernetes pod annotations                                         | `{}`                         |
| `deployment.replicaCount`   | Number of service replicas (pods)                                  | `1`                          |
| `deployment.strategy`       | Deployment strategy for updating pods                              | `RollingUpdate`              |
| `deployment.annotations`    | Kubernetes deployment annotations                                  | `{}`                         |
| `service.type`              | Kubernetes service type                                            | `ClusterIP`                  |
| `service.port`              | Kubernetes service port                                            | `4010`                       |
| `service.nodePort`          | Node port for *NodePort* service type                              | `nil`                        |
| `service.clusterIP`         | Cluster IP for *ClusterIP* service type                            | `nil`                        |
| `service.loadBalancerIP`    | Load balancer IP for *LoadBalancer* service type                   | `nil`                        |
| `service.externalIPs`       | Kubernetes external IPs that route to nodes if any                 | `[]`                         |
| `service.annotations`       | Kubernetes service annotations                                     | `{}`                         |
| `ingress.annotations`       | Kubernetes ingress annotations                                     | `{}`                         |
| `resources`                 | Resource requests and limits                                       | `{}`                         |
| `mongodb.enabled`           | Enable `mongodb` as a chart dependency                             | `true`                       |
| `mongodb.mongodbDatabase`   | MongoDB database name                                              | `sites`                      |
| `mongodb.mongodbUsername`   | MongoDB username                                                   | `site-service`               |
| `mongodb.extrenal.uri`      | External MongoDB URI                                               | `nil`                        |
| `mongodb.extrenal.username` | External MongoDB username                                          | `nil`                        |
| `mongodb.extrenal.secret`   | Kubernetes secret with `mongodb-password` key for external MongoDB | `nil`                        |
| `jaeger.enabled`            | Enable `jaeger-agent` sidecar for tracing                          | `false`                      |
| `jaeger.collector.address`  | Host and port for `jaeger-collector`                               | `nil`                        |
| `jaeger.agent.image`        | Docker image repository for `jaeger-agent`                         | `jaegertracing/jaeger-agent` |
| `jaeger.agent.tag`          | Docker image tag for `jaeger-agent`                                | `latest`                     |
| `jaeger.agent.pullPolicy`   | Docker image pull policy for `jaeger-agent`                        | `Always`                     |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example,

```bash
$ helm install --name my-release --set service.type=NodePort,service.nodePort=8080 repo/site
```

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example,

```bash
$ helm install --name my-release -f values.yaml repo/site
```

> **Tip**: You can use the default [values.yaml](values.yaml)
