#!/bin/bash

#
# USAGE:
#   ./bootstrap.sh int
#   ./bootstrap.sh manager
#   ./bootstrap.sh worker
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
  while [[ $# > 0 ]]; do
    case $1 in
      init|manager|worker) role="$1" ;;
    esac
    shift
  done

  role=${role:-""}
  whitelist_variable "node role" "init manager worker" "$role"
}


function init_swarm {
  addr=$(hostname -i | cut -d' ' -f2)

  docker swarm init --advertise-addr $addr &> /dev/null
  docker swarm join-token -q manager > /vagrant/swarm-manager-token
  docker swarm join-token -q worker > /vagrant/swarm-worker-token
}

function join_manager {
  addr=$(hostname -i | cut -d' ' -f2)
  token=$(cat /vagrant/swarm-manager-token)

  docker swarm join --token $token --advertise-addr $addr --listen-addr $addr:2377 manager-1:2377
}

function join_worker {
  addr=$(hostname -i | cut -d' ' -f2)
  token=$(cat /vagrant/swarm-worker-token)

  docker swarm join --token $token --advertise-addr $addr --listen-addr $addr:2377 manager-1:2377
}


ensure_command "docker"
process_args "$@"

case $role in
  init)     init_swarm    ;;
  manager)  join_manager  ;;
  worker)   join_worker   ;;
esac
