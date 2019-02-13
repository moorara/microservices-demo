ARG version

# BUILD STAGE
FROM golang:1.11.5-alpine as builder

RUN apk add --no-cache ca-certificates git && \
    go get -d github.com/caddyserver/builds && \
    go get -d github.com/mholt/caddy && \
    cd /go/src/github.com/mholt/caddy && \
    git checkout $version

WORKDIR /go/src/github.com/mholt/caddy/caddy

COPY plugins.list .
RUN while read plugin; do \
      sed -i "/This is where other plugins get plugged in (imported)/a \\\t_ \"$plugin\"" caddymain/run.go; \
    done < plugins.list && \
    go get -d ./...

RUN go run build.go && \
    cp ./caddy /


# FINAL STAGE
FROM alpine:3.8

COPY --from=builder /caddy /bin/caddy
RUN apk add --no-cache ca-certificates
WORKDIR /www

EXPOSE 80 443
ENTRYPOINT [ "caddy" ]

# USER nobody
# EXPOSE 8080 8443
# ENTRYPOINT [ "caddy", "-http-port", "8080", "-https-port", "8443" ]
