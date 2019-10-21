#!/bin/bash

set -euo pipefail


function install {
  go get -u github.com/golang/protobuf/proto
  go get -u github.com/golang/protobuf/protoc-gen-go
}

function generate {
  cd proto
  protoc --go_out=plugins=grpc:. *.proto
}


install
generate
