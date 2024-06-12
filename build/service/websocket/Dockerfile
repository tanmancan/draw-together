# Builds websocket app

ARG APP_NAMESPACE
ARG APP_VERSION

FROM gobase as buildcert

RUN addgroup -g 401 -S websocket && \
    adduser -u 401 -S websocket -G websocket
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
RUN chown 401 /misc/cert/cert.key

FROM scratch

COPY --from=buildcert /etc/passwd /etc/passwd
COPY --from=buildcert /misc/cert /misc/cert

USER websocket

WORKDIR /app

COPY --from=buildcert /app/websocket ./

ARG SERVICE_WEBSOCKET_PORT

EXPOSE ${SERVICE_WEBSOCKET_PORT}

ENTRYPOINT [ "./websocket" ]

