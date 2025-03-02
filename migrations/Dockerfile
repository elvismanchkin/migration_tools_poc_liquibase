FROM liquibase/liquibase:alpine AS builder

USER root
# Install PostgreSQL client for health check
RUN apk add --no-cache postgresql-client

# Create an entrypoint script with proper permissions
COPY entrypoint.sh /tmp/entrypoint.sh
RUN chmod +x /tmp/entrypoint.sh

FROM liquibase/liquibase:alpine

# Copy the PostgreSQL client from the builder stage
COPY --from=builder /usr/bin/pg_isready /usr/bin/pg_isready
COPY --from=builder /usr/lib/libpq.so* /usr/lib/
# Copy the entrypoint script with permissions already set
COPY --from=builder /tmp/entrypoint.sh /liquibase/entrypoint.sh

WORKDIR /liquibase

# Copy migrations folder to the container
COPY migrations/ /liquibase/changelog/migrations/
COPY master-changelog.yaml /liquibase/changelog/

# Set entrypoint
ENTRYPOINT ["/liquibase/entrypoint.sh"]