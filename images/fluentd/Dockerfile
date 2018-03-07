ARG fluentd_version

FROM fluent/fluentd:$fluentd_version
RUN apk add --update --virtual .build-deps sudo build-base ruby-dev && \
    sudo gem install fluent-plugin-elasticsearch && \
    sudo gem sources --clear-all && \
    apk del .build-deps && \
    rm -rf /var/cache/apk/* /home/fluent/.gem
