#!/bin/bash

#
# USAGE:
#   ./bootstrap.sh --addr 10.10.2.101 --name server-1 --count 3
#   ./bootstrap.sh --addr 10.10.2.102 --name server-2
#   ./bootstrap.sh -a 10.10.2.202 -n client-1
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
  while [[ $# > 1 ]]; do
    case $1 in
      -a|--addr)
      addr="$2"
      shift
      ;;
      -n|--name)
      name="$2"
      shift
      ;;
      -c|--count)
      count="$2"
      shift
      ;;
    esac
    shift
  done

  addr=${addr:-""}
  name=${name:-""}
  count=${count:-""}
}


function bootstrap {
  docker run -d \
    --net host \
    -v "/consul/data:/consul/data" \
    consul agent \
      -server \
      -bind=$addr \
      -bootstrap-expect=$count \
      -data-dir=/consul/data \
      -datacenter=local \
      -node=$name \
      -ui \
  &> /dev/null
}

function join_server {
  docker run -d \
    --net host \
    -v "/consul/data:/consul/data" \
    consul agent \
      -server \
      -bind=$addr \
      -retry-join=server-1 \
      -data-dir=/consul/data \
      -datacenter=local \
      -node=$name \
      -ui \
  &> /dev/null
}

function join_client {
  docker run -d \
    --net host \
    -v "/consul/data:/consul/data" \
    consul agent \
      -bind=$addr \
      -retry-join=server-1 \
      -data-dir=/consul/data \
      -datacenter=local \
      -node=$name \
      -ui \
  &> /dev/null
}


ensure_command "docker"
process_args "$@"

if [ "$name" == "server-1" ]; then
  bootstrap
elif [[ "$name" =~ ^server ]]; then
  join_server
elif [[ "$name" =~ ^client ]]; then
  join_client
fi
