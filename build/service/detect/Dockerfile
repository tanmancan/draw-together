# Builds detect app

ARG APP_NAMESPACE
ARG APP_VERSION

FROM gobase as buildcert

RUN addgroup -g 401 -S detect && \
    adduser -u 401 -S detect -G detect

FROM scratch

COPY --from=buildcert /etc/passwd /etc/passwd
COPY --from=buildcert /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER detect

COPY --chown=detect:detect --from=buildcert /tmp /tmp

WORKDIR /app

COPY --from=buildcert /app/detect ./

ENTRYPOINT [ "./detect" ]

