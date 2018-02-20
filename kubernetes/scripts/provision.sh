#!/bin/bash

set -euo pipefail


apt-get update
apt-get install -y vim jq

# Install rkt (https://coreos.com/rkt/docs/latest/distributions.html#deb-based)
gpg --list-keys &> /dev/null
gpg --recv-key 18AD5014C99EF7E3BA5F6CE950BDD3E0FC8A365E &> /dev/null
wget --quiet https://github.com/rkt/rkt/releases/download/v1.29.0/rkt_1.29.0-1_amd64.deb
wget --quiet https://github.com/rkt/rkt/releases/download/v1.29.0/rkt_1.29.0-1_amd64.deb.asc
gpg --verify rkt_1.29.0-1_amd64.deb.asc &> /dev/null
dpkg -i rkt_1.29.0-1_amd64.deb
groupadd rkt || true
usermod -aG rkt vagrant
