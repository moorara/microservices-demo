#!/bin/bash

#
# USAGE:
#   ./integration.sh
#

set -euo pipefail


red='\033[1;31m'
green='\033[1;32m'
yellow='\033[1;33m'
purple='\033[1;34m'
pink='\033[1;35m'
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

function run_react_client_tests {
  printf "\n${yellow}  RUNNING TESTS FOR ${blue}react-client${nocolor}\n"

  curl --fail http://localhost:4000/health &> /dev/null
  printf "${green}    [✓] GET /health ${nocolor}\n"

  curl --fail \
    -X GET \
    http://localhost:4000 &> /dev/null
  printf "${green}    [✓] GET / ${nocolor}\n"
}

function run_site_service_tests {
  printf "\n${yellow}  RUNNING TESTS FOR ${blue}site-service${nocolor}\n"

  siteName=$(get_random_string)

  curl --fail http://localhost:4010/health &> /dev/null
  printf "${green}    [✓] GET /health ${nocolor}\n"

  curl --fail http://localhost:4010/metrics &> /dev/null
  printf "${green}    [✓] GET /metrics ${nocolor}\n"

  curl --fail \
    -H "Content-Type: application/json" \
    -X POST -d "{ \"name\":\"$siteName\", \"location\":\"here\" }" \
    http://localhost:4010/v1/sites &> /dev/null
  printf "${green}    [✓] POST /v1/sites ${nocolor}\n"

  curl --fail \
    -H "Content-Type: application/json" \
    -X GET \
    http://localhost:4010/v1/sites &> /dev/null
  printf "${green}    [✓] GET /v1/sites ${nocolor}\n"
}

function run_sensor_service_tests {
  printf "\n${yellow}  RUNNING TESTS FOR ${blue}sensor-service${nocolor}\n"

  siteId=$(get_random_string)

  curl --fail http://localhost:4020/health &> /dev/null
  printf "${green}    [✓] GET /health ${nocolor}\n"

  curl --fail http://localhost:4020/metrics &> /dev/null
  printf "${green}    [✓] GET /metrics ${nocolor}\n"

  curl --fail \
    -H "Content-Type: application/json" \
    -X POST -d "{ \"siteId\":\"$siteId\", \"name\":\"temperature\", \"unit\":\"celsius\", \"minSafe\":-30, \"maxSafe\":30 }" \
    http://localhost:4020/v1/sensors &> /dev/null
  printf "${green}    [✓] POST /v1/sensors ${nocolor}\n"

  curl --fail \
    -H "Content-Type: application/json" \
    -X GET \
    http://localhost:4020/v1/sensors?siteId=$siteId &> /dev/null
  printf "${green}    [✓] GET /v1/sensors?siteId ${nocolor}\n"
}

function run_traefik_http_tests {
  printf "\n${yellow}  RUNNING INTEGRATION TESTS FOR ${pink}traefik (http)${nocolor}\n"

  siteName=$(get_random_string)
  siteId=$(get_random_string)

  curl --fail \
    -H "Host: traefik" \
    -H "Content-Type: application/json" \
    -X POST -d "{ \"name\":\"$siteName\", \"location\":\"here\" }" \
    http://localhost:1080/api/v1/sites &> /dev/null
  printf "${green}    [✓] POST /api/v1/sites ${nocolor}\n"

  curl --fail \
    -H "Host: traefik" \
    -H "Content-Type: application/json" \
    -X GET \
    http://localhost:1080/api/v1/sites &> /dev/null
  printf "${green}    [✓] GET /api/v1/sites ${nocolor}\n"

  curl --fail \
    -H "Host: traefik" \
    -H "Content-Type: application/json" \
    -X POST -d "{ \"siteId\":\"$siteId\", \"name\":\"temperature\", \"unit\":\"celsius\", \"minSafe\":-30, \"maxSafe\":30 }" \
    http://localhost:1080/api/v1/sensors &> /dev/null
  printf "${green}    [✓] POST /api/v1/sensors ${nocolor}\n"

  curl --fail \
    -H "Host: traefik" \
    -H "Content-Type: application/json" \
    -X GET \
    http://localhost:1080/api/v1/sensors?siteId=$siteId &> /dev/null
  printf "${green}    [✓] GET /api/v1/sensors?siteId ${nocolor}\n"
}

function run_traefik_https_tests {
  printf "\n${yellow}  RUNNING INTEGRATION TESTS FOR ${pink}traefik (https)${nocolor}\n"

  siteName=$(get_random_string)
  siteId=$(get_random_string)

  curl --fail \
    --cacert $certs_path/intermediate.ca.chain \
    -H "Host: traefik" \
    -H "Content-Type: application/json" \
    -X POST -d "{ \"name\":\"$siteName\", \"location\":\"here\" }" \
    https://localhost:1443/api/v1/sites &> /dev/null
  printf "${green}    [✓] POST /api/v1/sites ${nocolor}\n"

  curl --fail \
    --cacert $certs_path/intermediate.ca.chain \
    -H "Host: traefik" \
    -H "Content-Type: application/json" \
    -X GET \
    https://localhost:1443/api/v1/sites &> /dev/null
  printf "${green}    [✓] GET /api/v1/sites ${nocolor}\n"

  curl --fail \
    --cacert $certs_path/intermediate.ca.chain \
    -H "Host: traefik" \
    -H "Content-Type: application/json" \
    -X POST -d "{ \"siteId\":\"$siteId\", \"name\":\"temperature\", \"unit\":\"celsius\", \"minSafe\":-30, \"maxSafe\":30 }" \
    https://localhost:1443/api/v1/sensors &> /dev/null
  printf "${green}    [✓] POST /api/v1/sensors ${nocolor}\n"

  curl --fail \
    --cacert $certs_path/intermediate.ca.chain \
    -H "Host: traefik" \
    -H "Content-Type: application/json" \
    -X GET \
    https://localhost:1443/api/v1/sensors?siteId=$siteId &> /dev/null
  printf "${green}    [✓] GET /api/v1/sensors?siteId ${nocolor}\n"
}

function run_caddy_https_tests {
  printf "\n${yellow}  RUNNING INTEGRATION TESTS FOR ${purple}caddy (https)${nocolor}\n"

  siteName=$(get_random_string)
  siteId=$(get_random_string)

  curl --fail \
    --cacert $certs_path/intermediate.ca.chain \
    -X GET \
    https://localhost &> /dev/null
  printf "${green}    [✓] GET / ${nocolor}\n"

  curl --fail \
    --cacert $certs_path/intermediate.ca.chain \
    -H "Content-Type: application/json" \
    -X POST -d "{ \"name\":\"$siteName\", \"location\":\"here\" }" \
    https://localhost/api/v1/sites &> /dev/null
  printf "${green}    [✓] POST /api/v1/sites ${nocolor}\n"

  curl --fail \
    --cacert $certs_path/intermediate.ca.chain \
    -H "Content-Type: application/json" \
    -X GET \
    https://localhost/api/v1/sites &> /dev/null
  printf "${green}    [✓] GET /api/v1/sites ${nocolor}\n"

  curl --fail \
    --cacert $certs_path/intermediate.ca.chain \
    -H "Content-Type: application/json" \
    -X POST -d "{ \"siteId\":\"$siteId\", \"name\":\"temperature\", \"unit\":\"celsius\", \"minSafe\":-30, \"maxSafe\":30 }" \
    https://localhost/api/v1/sensors &> /dev/null
  printf "${green}    [✓] POST /api/v1/sensors ${nocolor}\n"

  curl --fail \
    --cacert $certs_path/intermediate.ca.chain \
    -H "Content-Type: application/json" \
    -X GET \
    https://localhost/api/v1/sensors?siteId=$siteId &> /dev/null
  printf "${green}    [✓] GET /api/v1/sensors?siteId ${nocolor}\n"
}


certs_path="../certs"

ensure_command "curl"
run_react_client_tests
run_site_service_tests
run_sensor_service_tests
run_traefik_http_tests
run_traefik_https_tests
run_caddy_https_tests
echo
