docker_tag ?= latest
docker_image ?= moorara/api-test
docker_dir := /tmp/workspace


clean:
	@ rm -rf *.log coverage

dep:
	@ dep ensure -update

run:
	@ go run cmd/main.go

build:
	@ ./scripts/build.sh ./cmd/main.go ./api-test

docker:
	@ cp -r ../../.git .git
	@ docker build -t $(docker_image):$(docker_tag) .
	@ rm -rf .git

test:
	@ go test -v -race ./...

coverage:
	@ ./scripts/test-unit-cover.sh

push:
	@ docker push $(docker_image):$(docker_tag)

save-images:
	@ mkdir -p $(docker_dir)
	@ docker image save -o $(docker_dir)/docker.tar $(docker_image):$(docker_tag)


.PHONY: clean
.PHONY: dep run build
.PHONY: docker
.PHONY: test coverage
.PHONY: push save-images
