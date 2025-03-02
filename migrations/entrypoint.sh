#!/bin/sh
set -e

echo "Waiting for database to be ready..."
until pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME; do
  echo "Database not ready yet, waiting..."
  sleep 2
done

echo "Running database migrations for $ENVIRONMENT environment..."

# Set up Liquibase properties
cat > /liquibase/liquibase.properties << EOF
changeLogFile=changelog/master-changelog.yaml
driver=org.postgresql.Driver
url=jdbc:postgresql://${DB_HOST}:${DB_PORT}/${DB_NAME}
username=${DB_USER}
password=${DB_PASSWORD}
defaultSchemaName=${DB_SCHEMA}
liquibase.hub.mode=off
# Include environment-specific context if needed
contexts=${ENVIRONMENT}
EOF

# Run Liquibase update
liquibase --defaultsFile=/liquibase/liquibase.properties update

echo "Migrations completed successfully."