# BUILD STAGE
FROM golang:1.10-alpine as builder
RUN apk add --no-cache ca-certificates make git
WORKDIR /go/src/github.com/moorara/microservices-demo/services/switch-service/
COPY . .
RUN make build && \
    cp switch-service /

# FINAL STAGE
FROM alpine:3.8
EXPOSE 4030 4031
HEALTHCHECK --interval=5s --timeout=3s --retries=3 CMD wget -q -O - http://localhost:4031/live || exit 1
RUN apk add --no-cache ca-certificates
COPY --from=builder /switch-service /usr/local/bin/
RUN chown -R nobody:nogroup /usr/local/bin/switch-service
USER nobody
CMD switch-service
