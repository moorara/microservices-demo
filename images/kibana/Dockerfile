ARG version

FROM docker.elastic.co/kibana/kibana-oss:$version
RUN ./bin/kibana-plugin install https://github.com/sivasamyk/logtrail/releases/download/v0.1.29/logtrail-6.4.0-0.1.29.zip
