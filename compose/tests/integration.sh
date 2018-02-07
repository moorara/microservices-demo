#!/bin/bash

#
# USAGE:
#   ./integration.sh
#

set -euo pipefail


red='\033[1;31m'
green='\033[1;32m'
yellow='\033[1;33m'
purple='\033[1;35m'
blue='\033[1;36m'
nocolor='\033[0m'


function get_random_string {
  echo $(cat /dev/random | LC_CTYPE=C tr -dc "[:alpha:]" | head -c 10)
}

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

function run_go_service_tests {
  printf "\n${yellow}  RUNNING TESTS FOR ${blue}go-service${nocolor}\n"

  linkId=$(get_random_string)

  curl --fail http://localhost:4010/health &> /dev/null
  printf "${green}    [✓] GET /health ${nocolor}\n"

  curl --fail http://localhost:4010/metrics &> /dev/null
  printf "${green}    [✓] GET /metrics ${nocolor}\n"

  curl --fail \
    -H "Content-Type: application/json" \
    -X POST -d "{ \"linkId\":\"$linkId\", \"stars\":4 }" \
    http://localhost:4010/v1/votes &> /dev/null
  printf "${green}    [✓] POST /v1/votes ${nocolor}\n"

  curl --fail \
    -H "Content-Type: application/json" \
    -X GET \
    http://localhost:4010/v1/votes?linkId=$linkId &> /dev/null
  printf "${green}    [✓] GET /v1/votes?linkId ${nocolor}\n"
}

function run_node_service_tests {
  printf "\n${yellow}  RUNNING TESTS FOR ${blue}node-service${nocolor}\n"

  linkURL=$(get_random_string)

  curl --fail http://localhost:4020/health &> /dev/null
  printf "${green}    [✓] GET /health ${nocolor}\n"

  curl --fail http://localhost:4020/metrics &> /dev/null
  printf "${green}    [✓] GET /metrics ${nocolor}\n"

  curl --fail \
    -H "Content-Type: application/json" \
    -X POST -d "{ \"url\":\"$linkURL\", \"title\":\"required\" }" \
    http://localhost:4020/v1/links &> /dev/null
  printf "${green}    [✓] POST /v1/links ${nocolor}\n"

  curl --fail \
    -H "Content-Type: application/json" \
    -X GET \
    http://localhost:4020/v1/links &> /dev/null
  printf "${green}    [✓] GET /v1/links ${nocolor}\n"
}

function run_traefik_tests {
  printf "\n${yellow}  RUNNING INTEGRATION TESTS FOR ${purple}traefik${nocolor}\n"

  linkId=$(get_random_string)
  curl --fail \
    -H "Host: traefik" \
    -H "Content-Type: application/json" \
    -X POST -d "{ \"linkId\":\"$linkId\", \"stars\":4 }" \
    http://localhost:1080/api/v1/votes &> /dev/null
  printf "${green}    [✓] POST /api/v1/votes ${nocolor}\n"

  curl --fail \
    -H "Host: traefik" \
    -H "Content-Type: application/json" \
    -X GET \
    http://localhost:1080/api/v1/votes?linkId=$linkId &> /dev/null
  printf "${green}    [✓] GET /api/v1/votes?linkId ${nocolor}\n"

  linkURL=$(get_random_string)
  curl --fail \
    -H "Host: traefik" \
    -H "Content-Type: application/json" \
    -X POST -d "{ \"url\":\"$linkURL\", \"title\":\"required\" }" \
    http://localhost:1080/api/v1/links &> /dev/null
  printf "${green}    [✓] POST /api/v1/links ${nocolor}\n"

  curl --fail \
    -H "Host: traefik" \
    -H "Content-Type: application/json" \
    -X GET \
    http://localhost:1080/api/v1/links &> /dev/null
  printf "${green}    [✓] GET /api/v1/links ${nocolor}\n"

  linkId=$(get_random_string)
  curl --fail \
    --cacert $certs_path/intermediate.ca.chain \
    -H "Host: traefik" \
    -H "Content-Type: application/json" \
    -X POST -d "{ \"linkId\":\"$linkId\", \"stars\":4 }" \
    https://localhost:1443/api/v1/votes &> /dev/null
  printf "${green}    [✓] POST /api/v1/votes (https) ${nocolor}\n"

  curl --fail \
    --cacert $certs_path/intermediate.ca.chain \
    -H "Host: traefik" \
    -H "Content-Type: application/json" \
    -X GET \
    https://localhost:1443/api/v1/votes?linkId=$linkId &> /dev/null
  printf "${green}    [✓] GET /api/v1/votes?linkId (https) ${nocolor}\n"

  linkURL=$(get_random_string)
  curl --fail \
    --cacert $certs_path/intermediate.ca.chain \
    -H "Host: traefik" \
    -H "Content-Type: application/json" \
    -X POST -d "{ \"url\":\"$linkURL\", \"title\":\"required\" }" \
    https://localhost:1443/api/v1/links &> /dev/null
  printf "${green}    [✓] POST /api/v1/links (https) ${nocolor}\n"

  curl --fail \
    --cacert $certs_path/intermediate.ca.chain \
    -H "Host: traefik" \
    -H "Content-Type: application/json" \
    -X GET \
    https://localhost:1443/api/v1/links &> /dev/null
  printf "${green}    [✓] GET /api/v1/links (https) ${nocolor}\n"
}


certs_path="../certs"

ensure_command "curl"
run_go_service_tests
run_node_service_tests
run_traefik_tests
echo
