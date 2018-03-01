#!/bin/bash

set -euo pipefail


apt-get update
apt-get install -y vim jq unzip

# HashiCorp PGP Key
gpg --import /vagrant/hashicorp.asc &> /dev/null

# Install Consul
{
  CONSUL_VERSION=1.0.6
  curl -sSL -O https://releases.hashicorp.com/consul/${CONSUL_VERSION}/consul_${CONSUL_VERSION}_linux_amd64.zip
  curl -sSL -O https://releases.hashicorp.com/consul/${CONSUL_VERSION}/consul_${CONSUL_VERSION}_SHA256SUMS.sig
  curl -sSL -O https://releases.hashicorp.com/consul/${CONSUL_VERSION}/consul_${CONSUL_VERSION}_SHA256SUMS
  gpg --verify consul_${CONSUL_VERSION}_SHA256SUMS.sig consul_${CONSUL_VERSION}_SHA256SUMS
  grep "consul_${CONSUL_VERSION}_linux_amd64.zip" consul_${CONSUL_VERSION}_SHA256SUMS | shasum -a 256
  unzip consul_${CONSUL_VERSION}_linux_amd64.zip -d /usr/local/bin/
  chmod 755 /usr/local/bin/consul
  rm -f consul_${CONSUL_VERSION}_*
} &> /dev/null

# Install Vault
{
  VAULT_VERSION=0.9.5
  curl -sSL -O https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_linux_amd64.zip
  curl -sSL -O https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_SHA256SUMS.sig
  curl -sSL -O https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_SHA256SUMS
  gpg --verify vault_${VAULT_VERSION}_SHA256SUMS.sig vault_${VAULT_VERSION}_SHA256SUMS
  grep "vault_${VAULT_VERSION}_linux_amd64.zip" vault_${VAULT_VERSION}_SHA256SUMS | shasum -a 256
  unzip vault_${VAULT_VERSION}_linux_amd64.zip -d /usr/local/bin/
  chmod 755 /usr/local/bin/vault
  rm -f vault_${VAULT_VERSION}_*
} &> /dev/null

# Install Nomad
{
  NOMAD_VERSION=0.7.1
  curl -sSL -O https://releases.hashicorp.com/nomad/${NOMAD_VERSION}/nomad_${NOMAD_VERSION}_linux_amd64.zip
  curl -sSL -O https://releases.hashicorp.com/nomad/${NOMAD_VERSION}/nomad_${NOMAD_VERSION}_SHA256SUMS.sig
  curl -sSL -O https://releases.hashicorp.com/nomad/${NOMAD_VERSION}/nomad_${NOMAD_VERSION}_SHA256SUMS
  gpg --verify nomad_${NOMAD_VERSION}_SHA256SUMS.sig nomad_${NOMAD_VERSION}_SHA256SUMS
  grep "nomad_${NOMAD_VERSION}_linux_amd64.zip" nomad_${NOMAD_VERSION}_SHA256SUMS | shasum -a 256
  unzip nomad_${NOMAD_VERSION}_linux_amd64.zip -d /usr/local/bin/
  chmod 755 /usr/local/bin/nomad
  rm -f nomad_${NOMAD_VERSION}_*
} &> /dev/null
