path := $(shell pwd)


up:
	@ mkdir -p vagrant && \
	  vagrant up

down:
	@ vagrant destroy --force

clean:
	@ rm -rf .vagrant vagrant

.PHONY: up down clean
