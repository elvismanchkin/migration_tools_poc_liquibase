#!/bin/sh
set -e

echo "Waiting for database to be ready..."
until pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME; do
  echo "Database not ready yet, waiting..."
  sleep 2
done

echo "Running database migrations for $ENVIRONMENT environment..."

FLYWAY_PLACEHOLDERS=""
if [ "$ENVIRONMENT" = "dev" ]; then
  FLYWAY_PLACEHOLDERS="-placeholders.environment=dev"
elif [ "$ENVIRONMENT" = "test" ]; then
  FLYWAY_PLACEHOLDERS="-placeholders.environment=test"
elif [ "$ENVIRONMENT" = "prod" ]; then
  FLYWAY_PLACEHOLDERS="-placeholders.environment=prod"
fi

flyway \
  -url=jdbc:postgresql://${DB_HOST}:${DB_PORT}/${DB_NAME} \
  -user=${DB_USER} \
  -password=${DB_PASSWORD} \
  -defaultSchema=${DB_SCHEMA} \
  ${FLYWAY_PLACEHOLDERS} \
  migrate

echo "Migrations completed successfully."
