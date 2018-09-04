#!/bin/bash

#
# EXAMPLES:
#   ./connect.sh manager-1
#   ./connect.sh worker-1
#

set -euo pipefail

path=$(dirname $0)
source "$path/functions.sh"


function process_args {
  if [[ $# == 0 ]] || [[ "$1" == "" ]]; then
    printf "${red}swarm node is not specified.${nocolor}\n"
    exit 1
  fi

  node="$1"
}

function cleanup {
  close_docker_tunnel "$node" || true
}


ensure_command "vagrant"
process_args "$@"

trap cleanup EXIT
write_ssh_config "$node"
open_docker_tunnel "$node"
PS1="${green}swarm-$node ${blue}> ${nocolor}" bash
