# Builds nginx image used to serve UI application

ARG APP_NAMESPACE
ARG APP_VERSION
ARG ALPINE_VERSION
ARG NGINX_VERSION

FROM uibase as uisetup

FROM nginx:${NGINX_VERSION}-alpine${ALPINE_VERSION} as ui

RUN apk update && apk upgrade && \
    apk add openssl

RUN mkdir -p /misc/cert
RUN openssl req -x509 \
    -newkey rsa:4096 \
    -sha256 \
    -days 30 \
    -nodes \
    -keyout /misc/cert/cert.key \
    -out /misc/cert/cert.crt \
    -subj "/CN=localhost" \
    -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"
RUN chown 101 /misc/cert/cert.key

COPY --from=uisetup /app/dist /usr/share/nginx/html
COPY build/service/ui/conf/default.conf /etc/nginx/conf.d/

ARG SERVICE_UI_PORT

EXPOSE ${SERVICE_UI_PORT}