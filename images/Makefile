fluentd_version ?= v1.2.2
kibana_version ?= 6.3.1


all: fluentd kibana

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


.PHONY: all fluentd kibana
