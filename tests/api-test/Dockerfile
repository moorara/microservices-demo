# BUILD STAGE
FROM golang:1.11-alpine as builder
RUN apk add --no-cache make git
WORKDIR /go/src/github.com/moorara/microservices-demo/tests/api-test/
COPY . .
RUN make build && \
    cp api-test /

# FINAL STAGE
FROM alpine:3.8
RUN apk add --no-cache ca-certificates
COPY --from=builder /api-test /usr/local/bin/
RUN chown -R nobody:nogroup /usr/local/bin/api-test
USER nobody
CMD api-test
