#!/bin/bash

#
# USAGE:
#   ./seed-data.sh
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

function pull_images {
  printf "${nocolor}<======================================== Pull Docker Images ========================================>\n"

  docker pull arangodb/arangodb
  docker pull mongo
  docker pull postgres
  docker pull cockroachdb/cockroach

  printf "${nocolor}\n"
}

function seed_arango {
  printf "${green}<======================================== ArangoDB ========================================>\n"

  # Create database
  # switches database is created at this point by switch-service
  # echo 'db._createDatabase("switches");' | arangosh --server.endpoint tcp://arango:8529 --server.password pass

  docker run \
    --network compose_local --link arango \
    --volume "$(pwd)/data/arango/switches.json:/data/switches.json" \
    arangodb/arangodb \
      arangoimp \
        --server.endpoint tcp://arango:8529 \
        --server.username root \
        --server.password pass \
        --server.database switches \
        --collection switches \
        --create-collection true \
        --file /data/switches.json

  printf "${nocolor}\n"
}

function seed_mongo {
  printf "${yellow}<======================================== MongoDB ========================================>\n"

  docker run \
    --network compose_local --link mongo \
    --volume "$(pwd)/data/mongo/sites.json:/data/sites.json" \
    mongo \
      mongoimport \
        --host mongo \
        --port 27017 \
        --db sites \
        --collection sites \
        /data/sites.json

  printf "${nocolor}\n"
}

function seed_postgres {
  printf "${purple}<======================================== PostgreSQL ========================================>\n"

  chmod 600 ./data/postgres/pgpass

  docker run \
    --network compose_local --link postgres \
    --volume "$(pwd)/data/postgres/pgpass:/root/.pgpass" \
    --volume "$(pwd)/data/postgres/sensors.sql:/data/sensors.sql" \
    postgres \
      psql \
        --echo-all \
        --host postgres \
        --port 5432 \
        --username root \
        --dbname sensors \
        --file /data/sensors.sql

  printf "${nocolor}\n"
}

function seed_cockroach {
  printf "${blue}<======================================== CockroachDB ========================================>\n"

  docker run \
    --network compose_local --link cockroach \
    --env COCKROACH_HOST=cockroach \
    --volume "$(pwd)/data/cockroach/assets.sql:/data/assets.sql" \
    cockroachdb/cockroach \
      shell -c "/cockroach/cockroach.sh sql --insecure --database assets < /data/assets.sql"

  printf "${nocolor}\n"
}


ensure_command docker
echo
pull_images
seed_arango
seed_mongo
seed_postgres
seed_cockroach
