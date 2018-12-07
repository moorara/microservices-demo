#!/bin/sh

#
# USAGE:
#   ./build.sh -m main.go -b app
#   ./build.sh --main main.go --binary app
#

set -euo pipefail


process_args() {
  while [ $# -gt 1 ]; do
    key=$1
    case $key in
      -m|--main)
      main=$2
      shift
      ;;
      -b|--binary)
      binary=$2
      shift
      ;;
    esac
    shift
  done

  main=${main:-"main.go"}
  binary=${binary:-"build/app"}
}

build_binary() {
  go build -o "$binary" "$main"
}


process_args "$@"
build_binary
