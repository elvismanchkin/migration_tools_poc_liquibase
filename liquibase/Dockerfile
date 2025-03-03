FROM liquibase/liquibase:alpine AS builder

USER root
RUN apk add --no-cache postgresql-client

COPY entrypoint.sh /tmp/entrypoint.sh
RUN chmod +x /tmp/entrypoint.sh

FROM liquibase/liquibase:alpine

COPY --from=builder /usr/bin/pg_isready /usr/bin/pg_isready
COPY --from=builder /usr/lib/libpq.so* /usr/lib/
COPY --from=builder /tmp/entrypoint.sh /liquibase/entrypoint.sh

WORKDIR /liquibase

COPY migrations/ /liquibase/changelog/migrations/
COPY master-changelog.yaml /liquibase/changelog/

ENTRYPOINT ["/liquibase/entrypoint.sh"]