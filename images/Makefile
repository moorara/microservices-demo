go_version ?= 1.11-alpine
node_version ?= 10-alpine
alpine_version ?= 3.8

circleci_go_image := moorara/circleci-go:$(go_version)
circleci_node_image := moorara/circleci-node:$(node_version)
circleci_tools_image := moorara/circleci-tools:alpine-$(alpine_version)

fluentd_version ?= v1.2.5
kibana_version ?= 6.4.0
caddy_version ?= v0.11.0
nats_exporter_version ?= v0.1.0


circleci: circleci-go circleci-node circleci-tools
push-circleci: push-circleci-go push-circleci-node push-circleci-tools

circleci-go:
	@ cd circleci-go && \
	  docker build --build-arg version=$(go_version) --tag $(circleci_go_image) .

push-circleci-go:
	@ cd circleci-go && \
	  docker push $(circleci_go_image)

circleci-node:
	@ cd circleci-node && \
	  docker build --build-arg version=$(node_version) --tag $(circleci_node_image) .

push-circleci-node:
	@ cd circleci-node && \
	  docker push $(circleci_node_image)

circleci-tools:
	@ cd circleci-tools && \
	  docker build --build-arg version=$(alpine_version) --tag $(circleci_tools_image) .

push-circleci-tools:
	@ cd circleci-tools && \
	  docker push $(circleci_tools_image)


all: caddy fluentd kibana nats-exporter
push-all: push-caddy push-fluentd push-kibana push-nats-exporter

caddy:
	@ cd caddy && \
	  docker build --build-arg version=$(caddy_version) --tag moorara/caddy:$(caddy_version) .

push-caddy:
	@ cd caddy && \
	  docker push moorara/caddy:$(caddy_version)

fluentd:
	@ cd fluentd && \
	  docker build --build-arg version=$(fluentd_version) --tag moorara/fluentd:$(fluentd_version) .

push-fluentd:
	@ cd fluentd && \
	  docker push moorara/fluentd:$(fluentd_version)

kibana:
	@ cd kibana && \
	  docker build --build-arg version=$(kibana_version) --tag moorara/kibana:$(kibana_version) .

push-kibana:
	@ cd kibana && \
	  docker push moorara/kibana:$(kibana_version)

nats-exporter:
	@ cd nats-exporter && \
	  docker build --build-arg version=$(nats_exporter_version) --tag moorara/nats-exporter:$(nats_exporter_version) .

push-nats-exporter:
	@ cd nats-exporter && \
	  docker push moorara/nats-exporter:$(nats_exporter_version)


.PHONY: circleci push-circleci circleci-go push-circleci-go circleci-node push-circleci-node circleci-tools push-circleci-tools
.PHONY: all push-all caddy push-caddy fluentd push-fluentd kibana push-kibana nats-exporter push-nats-exporter
