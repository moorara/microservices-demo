# BUILD STAGE
FROM golang:1.14-alpine as builder
RUN apk add --no-cache make git
WORKDIR /repo
COPY . .
RUN CGO_ENABLED=0 ./scripts/build.sh --main main.go --binary switch

# FINAL STAGE
FROM alpine:3.12
EXPOSE 4030 4031
HEALTHCHECK --interval=5s --timeout=3s --retries=3 CMD wget -q -O - http://localhost:4031/live || exit 1
RUN apk add --no-cache ca-certificates
COPY --from=builder /repo/switch /usr/local/bin/
RUN chown -R nobody:nogroup /usr/local/bin/switch
USER nobody
CMD [ "switch" ]
