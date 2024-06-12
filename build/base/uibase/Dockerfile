# Builds React application

ARG ALPINE_VERSION
ARG NODE_VERSION

FROM node:${NODE_VERSION}-alpine${ALPINE_VERSION} as protobuild

RUN apk update && apk upgrade && \
    apk add make

WORKDIR /app

COPY package.json package-lock.json Makefile ./
COPY proto ./proto

RUN make buildprotots

FROM node:${NODE_VERSION}-alpine${ALPINE_VERSION} as uibuild

WORKDIR /app/ui

COPY ui/package.json ui/package-lock.json ./

RUN npm ci

COPY --from=protobuild /app/ui/src/proto-ts ./
COPY ui ./

RUN npm run build

FROM scratch as uibase

WORKDIR /app

COPY --from=uibuild /app/ui/dist/ ./dist/