# Generates go artifacts from protobuf messages

ARG ALPINE_VERSION
ARG GO_VERSION

# Builds the go protobuf output
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as protogo

RUN apk update && apk upgrade && \
    apk add curl && \
    apk add unzip && \
    apk add make

ARG GO_VERSION
ARG PROTO_VERSION
ARG PROTO_GEN_GO_VERSION
ARG PROTO_GEN_GO_GRPC_VERSION
ARG TARGETPLATFORM

# Install protoc (protobuf compiler)
ENV PB_REL="https://github.com/protocolbuffers/protobuf/releases"
ENV PB_VER=${PROTO_VERSION}

RUN if [ "$TARGETPLATFORM" = "linux/amd64" ];\
    then PB_ARCH=x86_64;\
    elif [ "$TARGETPLATFORM" = "linux/arm64" ];\
    then PB_ARCH=aarch_64;\
    else PB_ARCH=x86_64;fi;\
    echo downloading $PB_REL/download/v$PB_VER/protoc-$PB_VER-linux-$PB_ARCH.zip &&\
    curl -LO $PB_REL/download/v$PB_VER/protoc-$PB_VER-linux-$PB_ARCH.zip &&\
    unzip protoc-$PB_VER-linux-$PB_ARCH.zip -d /usr/local/

WORKDIR /app

ENV PROTO_GEN_GO_VERSION=${PROTO_GEN_GO_VERSION}
ENV PROTO_GEN_GO_GRPC_VERSION=${PROTO_GEN_GO_GRPC_VERSION}

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTO_GEN_GO_VERSION}
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${PROTO_GEN_GO_GRPC_VERSION}

COPY proto/ ./proto/.
COPY Makefile ./

RUN make buildprotogo