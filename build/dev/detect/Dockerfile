ARG ALPINE_VERSION
ARG GO_VERSION

# Compiles go binaries
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION}

RUN apk update && apk upgrade && \
    apk add make && \
    apk add openssl

RUN addgroup -g 401 -S rpc && \
    adduser -u 401 -S rpc -G rpc

USER rpc

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go install github.com/cosmtrek/air@latest

COPY cmd/detect/ ./cmd/detect
COPY internal/ ./internal/
COPY Makefile .air.toml ./

ARG GITBRANCH
ARG GITCOMMIT
ARG GITSHORT
ARG TARGETPLATFORM

ENV GITBRANCH ${GITBRANCH}
ENV GITCOMMIT ${GITCOMMIT}
ENV GITSHORT ${GITSHORT}
ENV TARGETPLATFORM ${TARGETPLATFORM}

ENV TARGET detect
ENV BUILD_TARGET ./${TARGET}

RUN ln -sf ${TARGET} main 

ENTRYPOINT [ "air" ]