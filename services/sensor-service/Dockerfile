# BUILD STAGE
FROM golang:1.9-alpine as builder
RUN apk add --no-cache ca-certificates make
WORKDIR /go/src/github.com/moorara/microservices-demo/services/sensor-service/
COPY . .
RUN make build && \
    cp ./sensor-service /

# FINAL STAGE
FROM alpine:3.7
EXPOSE 4020
HEALTHCHECK --interval=5s --timeout=3s --retries=3 CMD wget -q -O - http://localhost:4020/health || exit 1
RUN apk add --no-cache ca-certificates
COPY --from=builder /sensor-service /usr/local/bin/
RUN chown -R nobody:nogroup /usr/local/bin/sensor-service && chmod +x /usr/local/bin/sensor-service
USER nobody
CMD sensor-service
