# Generates go binaries for the application
FROM protobase

RUN apk update && apk upgrade && \
    apk add make && \
    apk add openssl

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY internal/ ./internal
COPY cmd/rpc/ ./cmd/rpc
COPY cmd/websocket/ ./cmd/websocket
COPY cmd/detect/ ./cmd/detect
COPY Makefile ./

ARG GITBRANCH
ARG GITCOMMIT
ARG GITSHORT
ARG TARGETPLATFORM

ENV TARGET=rpc
ENV BUILD_TARGET=./${TARGET}
RUN make buildgobin

ENV TARGET=websocket
ENV BUILD_TARGET=./${TARGET}
RUN make buildgobin

ENV TARGET=detect
ENV BUILD_TARGET=./${TARGET}
RUN make buildgobin