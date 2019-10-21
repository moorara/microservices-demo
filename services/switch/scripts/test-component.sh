#!/bin/bash

set -euo pipefail


function process_args {
  while [[ $# > 0 ]]; do
    key=$1
    case $key in
      -d|--docker)
      docker=true
      ;;
      -v|--version)
      version=$2
      shift
      ;;
    esac
    shift
  done

  docker=${docker:-false}
  version=${version:-"latest"}
}

function test_component {
  export COMPONENT_TEST=true
  go test -v ./test/component
}

function test_component_docker {
  export VERSION=$version
  docker-compose run component-test
	docker container logs switch-service | grep '^{' | jq . > component-tests.log
	docker-compose down
}


process_args "$@"

case $docker in
  false) test_component ;;
  true)  test_component_docker ;;
esac
