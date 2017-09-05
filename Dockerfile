FROM alpine

# Misc
LABEL Description="Protologbeat Docker image based on Alpine" Vendor="Alain Lefebvre" 
MAINTAINER Alain Lefebvre <hartfordfive@gmail.com>

ARG VERSION
ENV VERSION=$VERSION

RUN set -ex ;\
    # Ensure protologbeat user exists
    addgroup -S protologbeat && adduser -S -G protologbeat protologbeat ;\
    # Install dependencies
    apk --no-cache add gettext libc6-compat curl ;\
    # Hotfix for libc compat
    ln -s /lib /lib64 ;\
    cd /tmp ;\
    mkdir -p /opt/protologbeat/conf ;\
    mkdir -p /opt/protologbeat/ssl ;\
    curl -L https://github.com/hartfordfive/protologbeat/releases/download/${VERSION}/protologbeat-${VERSION}-linux-x86_64.tar.gz --output protologbeat-${VERSION}-linux-x86_64.tar.gz ;\
    tar -xvzf protologbeat-${VERSION}-linux-x86_64.tar.gz ;\
    mv /tmp/protologbeat-${VERSION}-linux-x86_64 /opt/protologbeat/protologbeat ;\
    rm -rf protologbeat-${VERSION}-linux-x86_64 && rm protologbeat-${VERSION}-linux-x86_64.tar.gz ;\
    # Fix permissions
    chown -R protologbeat:protologbeat /opt/protologbeat ;\
    chmod 750 /opt/protologbeat ;\
    chmod 700 /opt/protologbeat/ssl

ENV PATH=/opt/protologbeat:$PATH

COPY protologbeat-docker.yml /opt/protologbeat/conf/protologbeat.yml
COPY protologbeat.template-es2x.json /opt/protologbeat
COPY protologbeat.template.json /opt/protologbeat

WORKDIR /opt/protologbeat
USER protologbeat

CMD ["protologbeat", "-c", "/opt/protologbeat/conf/protologbeat.yml"]
