# Generates Envoy Front Proxy
ARG ALPINE_VERSION
ARG ENVOY_VERSION

FROM alpine:${ALPINE_VERSION} as buildtpl

RUN apk update && apk upgrade && \
    apk add make && \
    apk add openssl

WORKDIR /app

COPY ./build/service/envoy/conf ./

ARG ENVOY_PUBLIC_PORT
ARG TLS_CERT_PATH
ARG TLS_KEY_PATH
ARG SERVICE_GRPC
ARG SERVICE_GRPC_PORT
ARG SERVICE_WEBSOCKET
ARG SERVICE_WEBSOCKET_PORT
ARG SERVICE_UI
ARG SERVICE_UI_PORT

RUN mkdir -p /misc/cert
RUN openssl req -x509 \
    -newkey rsa:4096 \
    -sha256 \
    -days 30 \
    -nodes \
    -keyout ${TLS_KEY_PATH} \
    -out ${TLS_CERT_PATH} \
    -subj "/CN=localhost" \
    -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

RUN make

FROM envoyproxy/envoy:distroless-${ENVOY_VERSION} as envoy

COPY --chown=0:0 --chmod=664 --from=buildtpl /misc/cert /misc/cert
COPY --chown=0:0 --from=buildtpl /app/tmp /etc/envoy

ARG ENVOY_PUBLIC_PORT
ARG ENVOY_ADMIN_PORT

EXPOSE ${ENVOY_PUBLIC_PORT}
EXPOSE ${ENVOY_ADMIN_PORT}