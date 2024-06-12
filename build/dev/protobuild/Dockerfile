FROM protobase

RUN apk update && apk upgrade && \
    apk add make

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go install github.com/cosmtrek/air@latest

COPY internal/ ./internal/
COPY proto ./proto/
COPY Makefile ./
COPY ./build/dev/protobuild/conf/.air.proto.toml ./.air.toml

RUN touch main
RUN chmod +x main 
RUN echo 'echo "Done"' > main

ARG TARGETPLATFORM
ENV TARGETPLATFORM ${TARGETPLATFORM}

RUN make buildprotogo

ENTRYPOINT [ "air" ]