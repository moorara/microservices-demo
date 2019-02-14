# BUILD STAGE
FROM golang:1.11-alpine as builder
RUN apk add --no-cache make git
WORKDIR /workspace
COPY . .
RUN CGO_ENABLED=0 ./scripts/build.sh --main main.go --binary sensor-service

# FINAL STAGE
FROM alpine:3.9
EXPOSE 4020
HEALTHCHECK --interval=5s --timeout=3s --retries=3 CMD wget -q -O - http://localhost:4020/health || exit 1
RUN apk add --no-cache ca-certificates
COPY --from=builder /workspace/sensor-service /usr/local/bin/
RUN chown -R nobody:nogroup /usr/local/bin/sensor-service
USER nobody
CMD [ "sensor-service" ]
