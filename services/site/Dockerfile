# STAGE 1
FROM node:12.18-alpine as builder
RUN apk add --no-cache ca-certificates
WORKDIR /usr/src/app
COPY [ "package.json", "." ]
RUN npm install --production
COPY . .

# STAGE 2
FROM node:12.18-alpine
ENV NODE_ENV production
EXPOSE 4010
HEALTHCHECK --interval=5s --timeout=3s --retries=3 CMD wget -q -O - http://localhost:4010/health || exit 1
RUN apk add --no-cache ca-certificates
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app .
USER nobody
CMD [ "node", "server.js" ]
