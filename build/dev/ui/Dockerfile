# Run React/Vite app in development mode
ARG ALPINE_VERSION
ARG NODE_VERSION
ARG SERVICE_UI_PORT

FROM node:${NODE_VERSION}-alpine${ALPINE_VERSION} as protobuild

WORKDIR /app

COPY ./ui/package.json ./ui/package-lock.json ./

RUN npm ci

COPY ./ui/tsconfig.json ./ui/tsconfig.node.json ./ui/vite.config.ts ./ui/index.html ./ui/.eslintrc.cjs ./

ARG ENVOY_PUBLIC_PORT

ENV VITE_RPC_BASE_URL https://localhost:${ENVOY_PUBLIC_PORT}
ENV VITE_WS_PORT ${ENVOY_PUBLIC_PORT}
ENV VITE_WS_HOST wss://localhost

EXPOSE ${SERVICE_UI_PORT}

CMD [ "npm", "run", "dev" ]
