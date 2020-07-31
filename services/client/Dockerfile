# BUILD STAGE
FROM node:12.18-alpine as builder
RUN apk add --no-cache ca-certificates
WORKDIR /usr/src/app
COPY . .
RUN yarn install && yarn run build

# FINAL STAGE
FROM node:12.18-alpine
EXPOSE 4000
HEALTHCHECK --interval=5s --timeout=3s --retries=3 CMD wget -q -O - http://localhost:4000/health || exit 1
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/build/ ./build/
COPY --from=builder /usr/src/app/server/ ./server/
RUN cd server && yarn install
USER nobody
CMD [ "node", "server/index.js" ]
