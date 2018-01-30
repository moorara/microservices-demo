# STAGE 1
FROM node:9-alpine as builder
RUN apk add --no-cache ca-certificates
WORKDIR /usr/src/app
COPY [ "package.json", "." ]
RUN npm install --production
COPY . .

# STAGE 2
FROM node:9-alpine
EXPOSE 4020
HEALTHCHECK --interval=5s --timeout=3s --retries=3 CMD wget -q -O - http://localhost:4020/health || exit 1
RUN apk add --no-cache ca-certificates
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app .
USER nobody
CMD [ "node", "server.js" ]
