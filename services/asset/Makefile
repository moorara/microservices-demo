name := asset-service
version ?= latest
docker_image ?= moorara/$(name)
docker_test_image ?= $(name)-test


clean:
	@ rm -rf *.log coverage

run:
	@ go run main.go

build:
	@ ./scripts/build.sh --main main.go --binary $(name)

test:
	@ go test -race ./...

coverage:
	@ mkdir -p coverage
	@ go test -covermode=atomic -coverprofile=coverage/c.out ./...
	@ go tool cover -html=coverage/c.out -o coverage/coverage.html

docker:
	@ cp -r ../../.git .git
	@ docker build --tag $(docker_image):$(version) .
	@ rm -rf .git

docker-test:
	@ docker build --file Dockerfile.test --tag $(docker_test_image) .

push:
	@ docker image push $(docker_image):$(version)

save-docker:
	@ docker image save --output ./docker.tar $(docker_image):$(version)

load-docker:
	@ docker image load --input docker.tar

up:
	@ VERSION=$(version) docker-compose up -d $(name)

down:
	@ VERSION=$(version) docker-compose down

test-integration: # docker up
	@ ./scripts/test-integration.sh --version $(version)

test-integration-docker: # docker docker-test
	@ ./scripts/test-integration.sh --version $(version) --docker

test-component: # docker up
	@ ./scripts/test-component.sh --version $(version)

test-component-docker: # docker docker-test
	@ ./scripts/test-component.sh --version $(version) --docker


.PHONY: clean
.PHONY: run build
.PHONY: test coverage
.PHONY: docker docker-test push save-docker load-docker
.PHONY: up down test-integration test-integration-docker test-component test-component-docker
