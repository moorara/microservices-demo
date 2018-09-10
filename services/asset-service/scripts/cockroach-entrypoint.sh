#!/bin/sh

set -eu

{
  /cockroach/cockroach.sh user set --insecure service
  /cockroach/cockroach.sh sql --insecure -e 'CREATE DATABASE assets'
  /cockroach/cockroach.sh sql --insecure -e 'GRANT ALL ON DATABASE assets TO service'
} 1> /dev/null
