ARG REDIS_VERSION
ARG ALPINE_VERSION

FROM redis:${REDIS_VERSION}-alpine

ENTRYPOINT [ "redis-server" ]

CMD [ "--save", "60", "1", "--loglevel", "warning" ]