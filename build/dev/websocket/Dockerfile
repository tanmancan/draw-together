ARG ALPINE_VERSION
ARG GO_VERSION

# Compiles go binaries
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION}

RUN apk update && apk upgrade && \
    apk add make && \
    apk add openssl

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

USER websocket

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go install github.com/cosmtrek/air@latest

COPY cmd/websocket/ ./cmd/websocket
COPY internal/ ./internal/
COPY Makefile .air.toml ./

ARG GITBRANCH
ARG GITCOMMIT
ARG GITSHORT
ARG TARGETPLATFORM
ARG SERVICE_WEBSOCKET_PORT

ENV GITBRANCH ${GITBRANCH}
ENV GITCOMMIT ${GITCOMMIT}
ENV GITSHORT ${GITSHORT}
ENV TARGETPLATFORM ${TARGETPLATFORM}
ENV SERVICE_WEBSOCKET_PORT ${SERVICE_WEBSOCKET_PORT}

ENV TARGET websocket
ENV BUILD_TARGET ./${TARGET}

RUN ln -sf ${TARGET} main 

EXPOSE ${SERVICE_WEBSOCKET_PORT}

ENTRYPOINT [ "air" ]