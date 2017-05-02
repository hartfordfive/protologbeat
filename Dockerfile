FROM node:6-alpine

# Misc
LABEL Description="Protologbeat Docker image based on Alpine" Vendor="Alain Lefebvre" 
MAINTAINER Alain Lefebvre <hartfordfive@gmail.com>

#RUN apk update && \
#    apk upgrade && \
#    apk add curl

ARG version
ENV VERSION=$version

RUN set -ex ;\
    # Ensure kibana user exists
    addgroup -S protologbeat && adduser -S -G protologbeat protologbeat ;\
    # Install dependencies
    apk --no-cache add bash fontconfig gettext su-exec tini curl ;\
    # Fix permissions
    mkdir -p /opt/protologbeat/conf && mkdir -p /opt/protologbeat/ssl

RUN curl -Lso - https://github.com/hartfordfive/protologbeat/releases/download/${VERSION}/protologbeat-${VERSION}-linux-x86_64.tar.gz | \
      tar zxf - -C /tmp && \
      cp /tmp/protologbeat-${VERSION}-linux-x86_64 /opt/protologbeat/protologbeat
#    cp /tmp/protologbeat-0.1.0-linux-x86_64 /usr/share/protologbeat

ENV PATH=/opt/protologbeat:$PATH

COPY protologbeat-docker.yml /opt/protologbeat/conf/protologbeat.yml
COPY protologbeat.template-es2x.json /opt/protologbeat
COPY protologbeat.template.json /opt/protologbeat

RUN chown -R protologbeat:protologbeat /opt/protologbeat ;\
    chmod 750 /opt/protologbeat ;\
    chmod 700 /opt/protologbeat/ssl

WORKDIR /opt/protologbeat


USER protologbeat





CMD ["protologbeat", "-e", "-c", "/opt/protologbeat/conf/protologbeat.yml"]
