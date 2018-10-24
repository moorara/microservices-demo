ARG version

FROM fluent/fluentd:$version
RUN apk add --update --virtual .build-deps sudo build-base ruby-dev && \
    sudo gem install fluent-plugin-elasticsearch && \
    sudo gem sources --clear-all && \
    apk del .build-deps && \
    rm -rf /var/cache/apk/* /home/fluent/.gem
