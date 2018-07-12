# STAGE 1
FROM node:10-alpine as builder
RUN apk add --no-cache ca-certificates
WORKDIR /usr/src/app
COPY . .
RUN yarn install && \
    yarn run build:webpack

# STAGE 2
FROM node:10-alpine
EXPOSE 4000
HEALTHCHECK --interval=5s --timeout=3s --retries=3 CMD wget -q -O - http://localhost:4000/health || exit 1
WORKDIR /usr/src/app
COPY --from=builder /usr/src/app/public/ ./public/
COPY --from=builder /usr/src/app/node/prod.js ./node/
COPY --from=builder [ "/usr/src/app/package.json", "/usr/src/app/.babelrc", "./" ]
RUN yarn install --production
USER nobody
CMD [ "node", "node/prod.js" ]
