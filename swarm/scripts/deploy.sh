#!/bin/bash

#
# EXAMPLES:
#

set -euo pipefail

path=$(dirname $0)
source "$path/functions.sh"


function process_args {
  while [[ $# > 0 ]]; do
    key="$1"
    case $key in
      -s|--stack)
      stack="$2"
      shift
      ;;
      -n|--node)
      node="$2"
      shift
      ;;
    esac
    shift
  done

  stack=${stack:-""}
  node=${node-"manager-1"}
  compose_file="./docker-compose/docker-compose.$stack.yml"
}

function cleanup {
  close_docker_tunnel "$node" || true
}


trap cleanup EXIT

# Checks
process_args "$@"
ensure_variable "stack" "node"
ensure_command "docker"

# Open a tunnel to a manager
write_ssh_config "$node"
open_docker_tunnel "$node"

# Deploy the stack
docker stack services "$stack"
docker stack deploy --compose-file "$compose_file" --with-registry-auth "$stack"
docker stack services "$stack"
