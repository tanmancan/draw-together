ARG ALPINE_VERSION
ARG GO_VERSION

# Compiles go binaries
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION}

RUN apk update && apk upgrade && \
    apk add make && \
    apk add openssl

RUN addgroup -g 401 -S rpc && \
    adduser -u 401 -S rpc -G rpc
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

USER rpc

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go install github.com/cosmtrek/air@latest

COPY cmd/rpc/ ./cmd/rpc
COPY internal/ ./internal/
COPY Makefile .air.toml ./

ARG GITBRANCH
ARG GITCOMMIT
ARG GITSHORT
ARG TARGETPLATFORM
ARG SERVICE_GRPC_PORT

ENV GITBRANCH ${GITBRANCH}
ENV GITCOMMIT ${GITCOMMIT}
ENV GITSHORT ${GITSHORT}
ENV TARGETPLATFORM ${TARGETPLATFORM}
ENV SERVICE_GRPC_PORT ${SERVICE_GRPC_PORT}

ENV TARGET rpc
ENV BUILD_TARGET ./${TARGET}

RUN ln -sf ${TARGET} main 

EXPOSE ${SERVICE_GRPC_PORT}

ENTRYPOINT [ "air" ]