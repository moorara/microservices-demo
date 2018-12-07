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

function test_integration {
  export INTEGRATION_TEST=true
  go test -v ./test/integration
}

function test_integration_docker {
  export VERSION=$version
	docker-compose run integration-test
	docker-compose down
}


process_args "$@"

case $docker in
  false) test_integration ;;
  true)  test_integration_docker ;;
esac
