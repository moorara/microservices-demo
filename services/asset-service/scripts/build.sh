#!/bin/sh

set -euo pipefail


version=$(cat VERSION)
revision=$(git rev-parse --short HEAD)
branch=$(git rev-parse --abbrev-ref HEAD)
buildtime=$(date -u +%Y-%m-%dT%H:%M:%SZ)

version_flag="-X $(go list ./cmd/version).Version=$version"
revision_flag="-X $(go list ./cmd/version).Revision=$revision"
branch_flag="-X $(go list ./cmd/version).Branch=$branch"
buildtime_flag="-X $(go list ./cmd/version).BuildTime=$buildtime"

function build_binary {
  main_file=$1
  binary_file=$2

  go build \
    -o $binary_file \
    -ldflags "$version_flag $revision_flag $branch_flag $buildtime_flag" \
    $main_file
}


build_binary $1 $2
