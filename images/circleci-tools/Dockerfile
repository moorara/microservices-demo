ARG version

FROM alpine:$version
RUN apk add --no-cache \
      ca-certificates openssh openssh-client \
      bash make git wget curl jq tar gzip gnupg

# Install docker and docker-compose
RUN curl -fsSL https://raw.githubusercontent.com/moorara/scripts/master/docker/install-docker.sh | bash && \
    curl -fsSL https://raw.githubusercontent.com/moorara/scripts/master/docker/install-compose.sh | bash

# Install kubectl and helm
RUN curl -fsSL https://raw.githubusercontent.com/moorara/scripts/master/kubernetes/install-kubectl.sh | bash && \
    curl -fsSL https://raw.githubusercontent.com/moorara/scripts/master/kubernetes/install-helm.sh | bash
