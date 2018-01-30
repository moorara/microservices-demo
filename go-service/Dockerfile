# BUILD STAGE
FROM golang:1.9-alpine as builder
RUN apk add --no-cache ca-certificates make
WORKDIR /go/src/github.com/moorara/toys/microservices/go-service/
COPY . .
RUN make build && \
    cp ./go-service /

# FINAL STAGE
FROM alpine:3.7
EXPOSE 4010
HEALTHCHECK --interval=5s --timeout=3s --retries=3 CMD wget -q -O - http://localhost:4010/health || exit 1
RUN apk add --no-cache ca-certificates
COPY --from=builder /go-service /usr/local/bin/
RUN chown -R nobody:nogroup /usr/local/bin/go-service && chmod +x /usr/local/bin/go-service
USER nobody
CMD go-service
