ARG kibana_version

FROM docker.elastic.co/kibana/kibana-oss:$kibana_version
RUN ./bin/kibana-plugin install https://github.com/sivasamyk/logtrail/releases/download/v0.1.27/logtrail-6.2.2-0.1.27.zip
