fluentd_version ?= v1.2.5
kibana_version ?= 6.4.0
caddy_version ?= v0.11.0
nats_exporter_version ?= v0.1.0


all: fluentd kibana caddy nats-exporter

fluentd:
	@ cd fluentd && \
	  docker build \
	    --build-arg fluentd_version=$(fluentd_version) \
	    -t moorara/fluentd:$(fluentd_version) .

kibana:
	@ cd kibana && \
	  docker build \
	    --build-arg kibana_version=$(kibana_version) \
	    -t moorara/kibana:$(kibana_version) .

caddy:
	@ cd caddy && \
	  docker build \
	    --build-arg caddy_version=$(caddy_version) \
	    -t moorara/caddy:$(caddy_version) .

nats-exporter:
	@ cd nats-exporter && \
	  docker build \
		  --build-arg nats_exporter_version=$(nats_exporter_version) \
	    -t moorara/nats-exporter:$(nats_exporter_version) .


.PHONY: all fluentd kibana caddy nats-exporter
