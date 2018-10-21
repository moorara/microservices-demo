docker_tag ?= latest
docker_image ?= moorara/graphql-server
docker_test_image ?= graphql-server-test
docker_dir := /tmp/workspace


clean:
	@ rm -rf coverage .nyc_output component-tests.log

docker:
	@ docker build -t $(docker_image):$(docker_tag) .

docker-test:
	@ docker build \
		--file Dockerfile.test \
		--build-arg version=$(docker_tag) \
		--tag $(docker_test_image) .

up:
	@ VERSION=$(docker_tag) \
	  docker-compose up -d graphql-server

down:
	@ VERSION=$(docker_tag) \
	  docker-compose down

test-docker: # docker docker-test
	@ docker run --rm $(docker_test_image) npm run test

test-component-docker: # docker docker-test
	@ VERSION=$(docker_tag) \
	  docker-compose run component-test && \
		docker container logs graphql-server | grep '^{' | jq . > component-tests.log && \
		docker-compose down

push:
	@ docker push $(docker_image):$(docker_tag)

save-images:
	@ mkdir -p $(docker_dir)
	@ docker image save -o $(docker_dir)/docker.tar $(docker_image):$(docker_tag)
	@ docker image save -o $(docker_dir)/docker-test.tar $(docker_test_image)


.PHONY: clean
.PHONY: docker docker-test
.PHONY: up down
.PHONY: test-docker test-component-docker
.PHONY: push save-images
