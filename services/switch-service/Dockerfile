# BUILD STAGE
FROM golang:1.10-alpine as builder
RUN apk add --no-cache ca-certificates make git
WORKDIR /go/src/github.com/moorara/microservices-demo/services/switch-service/
COPY . .
RUN make build && \
    cp switch-service /

# FINAL STAGE
FROM alpine:3.8
EXPOSE 4030
RUN apk add --no-cache ca-certificates
COPY --from=builder /switch-service /usr/local/bin/
RUN chown -R nobody:nogroup /usr/local/bin/switch-service
USER nobody
CMD switch-service
