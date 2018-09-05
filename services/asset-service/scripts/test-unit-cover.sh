#!/bin/bash

set -euo pipefail


coverpath="coverage"
coverfile="coverage.out"
go_packages=$(go list ./... | grep -v //)

function prepare {
  mkdir -p $coverpath
  echo "mode: atomic" > $coverpath/$coverfile
}

function test_package {
  package=$1
  go test -covermode=atomic -coverprofile=$coverfile $package
  tail -n +2 $coverfile >> $coverpath/$coverfile
  rm $coverfile
}

function generate_html {
  go tool cover -html=$coverpath/$coverfile -o $coverpath/index.html
  rm $coverpath/$coverfile
}


prepare
for pkg in $go_packages; do
  test_package $pkg
done
generate_html
