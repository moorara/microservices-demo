fluentd_version ?= v1.3.1
caddy_version ?= v0.11.1
nats_exporter_version ?= v0.1.0


all: caddy fluentd nats-exporter
push-all: push-caddy push-fluentd push-nats-exporter

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

nats-exporter:
	@ cd nats-exporter && \
	  docker build --build-arg version=$(nats_exporter_version) --tag moorara/nats-exporter:$(nats_exporter_version) .

push-nats-exporter:
	@ cd nats-exporter && \
	  docker push moorara/nats-exporter:$(nats_exporter_version)


.PHONY: all push-all caddy push-caddy fluentd push-fluentd nats-exporter push-nats-exporter
