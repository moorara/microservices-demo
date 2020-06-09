FROM fluent/fluentd:v1.11

USER root
RUN apk add --no-cache --update --virtual .build-deps \
    sudo build-base ruby-dev && \
    fluent-gem install fluent-plugin-elasticsearch && \
    fluent-gem install fluent-plugin-prometheus && \
    sudo gem sources --clear-all && \
    apk del .build-deps && \
    rm -rf /tmp/* /var/tmp/* /usr/lib/ruby/gems/*/cache/*.gem

COPY fluent.conf /fluentd/etc/
USER fluent
