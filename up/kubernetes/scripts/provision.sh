#!/bin/bash

set -euo pipefail


apt-get update
apt-get install -y vim jq

# CoreOS PGP Key
gpg --import /vagrant/app-signing-pubkey.gpg &> /dev/null

# Install rkt
{
  RKT_VERSION=1.29.0
  curl -sSL -O https://github.com/rkt/rkt/releases/download/v${RKT_VERSION}/rkt_${RKT_VERSION}-1_amd64.deb
  curl -sSL -O https://github.com/rkt/rkt/releases/download/v${RKT_VERSION}/rkt_${RKT_VERSION}-1_amd64.deb.asc
  gpg --verify rkt_${RKT_VERSION}-1_amd64.deb.asc rkt_${RKT_VERSION}-1_amd64.deb
  dpkg -i rkt_1.29.0-1_amd64.deb
  groupadd rkt || true && usermod -aG rkt vagrant
  rm -f rkt_${RKT_VERSION}-*
} &> /dev/null

# Install etcd
{
  ETCD_VERSION=3.3.1
  curl -sSL -O https://github.com/coreos/etcd/releases/download/v${ETCD_VERSION}/etcd-v${ETCD_VERSION}-linux-amd64.tar.gz
  curl -sSL -O https://github.com/coreos/etcd/releases/download/v${ETCD_VERSION}/etcd-v${ETCD_VERSION}-linux-amd64.aci.asc
  curl -sSL -O https://github.com/coreos/etcd/releases/download/v${ETCD_VERSION}/etcd-v${ETCD_VERSION}-linux-amd64.aci
  gpg --verify etcd-v${ETCD_VERSION}-linux-amd64.aci.asc etcd-v${ETCD_VERSION}-linux-amd64.aci
  tar -xzvf etcd-v${ETCD_VERSION}-linux-amd64.tar.gz
  mv etcd-v${ETCD_VERSION}-linux-amd64/etcd /usr/local/bin/
  mv etcd-v${ETCD_VERSION}-linux-amd64/etcdctl /usr/local/bin/
  rm -rf etcd-v${ETCD_VERSION}-*
} &> /dev/null
