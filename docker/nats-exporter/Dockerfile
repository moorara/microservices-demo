ARG version=v0.6.0

# BUILD STAGE
FROM golang:1.14.6-alpine AS builder
RUN apk add --no-cache git && \
    go get -d github.com/nats-io/prometheus-nats-exporter && \
    cd /go/src/github.com/nats-io/prometheus-nats-exporter && \
    git checkout $version && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -tags netgo -installsuffix netgo -ldflags "-s -w"

# FINAL STAGE
FROM alpine:3.12
COPY --from=builder /go/src/github.com/nats-io/prometheus-nats-exporter/prometheus-nats-exporter /usr/local/bin/
RUN chown -R nobody:nogroup /usr/local/bin/prometheus-nats-exporter && chmod 755 /usr/local/bin/prometheus-nats-exporter
USER nobody
EXPOSE 7777
ENTRYPOINT [ "prometheus-nats-exporter" ]
