path := $(shell pwd)

red := \033[1;31m
green := \033[1;32m
yellow := \033[1;33m
purple := \033[1;35m
blue := \033[1;36m
nocolor := \033[0m

define print_title
	printf "$(yellow)\n=============================[ $(1) ]==============================\n\n$(nocolor)"
endef


ci-test-go-service:
	@ cd microservices/go-service && \
	  $(call print_title,RUNNING UNIT TESTS)       &&  make test            && \
	  $(call print_title,BUILDING DOCKER IMAGES)   &&  make docker up       && \
	  $(call print_title,RUNNING COMPONENT TESTS)  &&  make test-component  && \
	  $(call print_title,DONE)

ci-test-node-service:
	@ cd microservices/node-service && \
	  $(call print_title,INSTALLING DEPENDENCIES)  &&  yarn                        && \
	  $(call print_title,RUNNING NODE SECURITY)    &&  yarn run nsp                && \
	  $(call print_title,RUNNING STANDARD LINTER)  &&  yarn run lint               && \
	  $(call print_title,RUNNING UNIT TESTS)       &&  yarn run test               && \
	  $(call print_title,BUILDING DOCKER IMAGES)   &&  make docker docker-test up  && \
	  $(call print_title,RUNNING COMPONENT TESTS)  &&  yarn run test:component     && \
	  $(call print_title,DONE)


.PHONY: ci-test-go-service
.PHONY: ci-test-node-service
