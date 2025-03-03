FROM flyway/flyway:9.8-alpine

USER root
RUN apk add --no-cache postgresql-client

COPY entrypoint.sh /flyway/entrypoint.sh
RUN chmod +x /flyway/entrypoint.sh

COPY sql/ /flyway/sql/

ENTRYPOINT ["/flyway/entrypoint.sh"]
