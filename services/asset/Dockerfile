# BUILD STAGE
FROM golang:1.14-alpine as builder
RUN apk add --no-cache make git
WORKDIR /repo
COPY . .
RUN CGO_ENABLED=0 ./scripts/build.sh --main main.go --binary asset

# FINAL STAGE
FROM alpine:3.12
EXPOSE 4040
HEALTHCHECK --interval=5s --timeout=3s --retries=3 CMD wget -q -O - http://localhost:4040/liveness || exit 1
RUN apk add --no-cache ca-certificates
COPY --from=builder /repo/asset /usr/local/bin/
RUN chown -R nobody:nogroup /usr/local/bin/asset
USER nobody
CMD [ "asset" ]
