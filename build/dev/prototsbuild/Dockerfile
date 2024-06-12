
ARG ALPINE_VERSION
ARG NODE_VERSION

FROM node:${NODE_VERSION}-alpine${ALPINE_VERSION} as prototsbuild

RUN apk update && apk upgrade && \
    apk add make

WORKDIR /app

RUN npm install -g nodemon

COPY build/dev/prototsbuild/conf/nodemon.json .
COPY package.json package-lock.json Makefile ./
COPY proto ./proto

RUN npm ci

RUN mkdir -p ui/src/proto-ts

ENTRYPOINT [ "nodemon", "--config", "nodemon.json" ] 