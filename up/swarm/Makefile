path := $(shell pwd)


up:
	@ vagrant up

down:
	@ vagrant destroy --force

clean:
	@ rm -rf .vagrant vagrant/hosts vagrant/swarm-*-token

.PHONY: up down clean
