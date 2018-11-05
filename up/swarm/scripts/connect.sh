#!/bin/bash

#
# EXAMPLES:
#   ./connect.sh manager-1
#   ./connect.sh worker-1
#

set -euo pipefail


red='\033[1;31m'
green='\033[1;32m'
yellow='\033[1;33m'
purple='\033[1;35m'
blue='\033[1;36m'
nocolor='\033[0m'


function ensure_command {
  for cmd in $@; do
    which $cmd 1> /dev/null || (
      printf "${red}$cmd not available!${nocolor}\n"
      exit 1
    )
  done
}

function ensure_env_var {
  for var in $@; do
    if [ "${!var}" == "" ]; then
      printf "${red}$var is not set.${nocolor}\n"
      exit 1
    fi
  done
}

function whitelist_variable {
  if [[ ! $2 =~ (^|[[:space:]])$3($|[[:space:]]) ]]; then
    printf "${red}Invalid $1 $3${nocolor}\n"
    exit 1
  fi
}

function process_args {
  if [[ $# == 0 ]] || [[ "$1" == "" ]]; then
    printf "${red}swarm node is not specified.${nocolor}\n"
    exit 1
  fi
  node="$1"

  ssh_control_path=$(mktemp -u)
  ssh_config_file="$node-config"
}


function write_ssh_config {
  vagrant ssh-config $node > $ssh_config_file
  echo "  ControlPath $ssh_control_path" >> $ssh_config_file
}

function open_docker_tunnel {
  export DOCKER_HOST=127.0.0.1:23770
  ssh -F $ssh_config_file -M -fnNT -L $DOCKER_HOST:/var/run/docker.sock $node
}

function close_docker_tunnel {
  unset DOCKER_HOST
  ssh -F $ssh_config_file -O exit $node > /dev/null 2>&1 || true
  rm -f $ssh_config_file $ssh_control_path > /dev/null 2>&1 || true
}

function cleanup {
  close_docker_tunnel || true
}


ensure_command "vagrant"
process_args "$@"

trap cleanup EXIT
write_ssh_config
open_docker_tunnel
PS1="${green}swarm-$node ${blue}> ${nocolor}" bash
