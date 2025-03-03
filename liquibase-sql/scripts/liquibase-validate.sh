#!/bin/bash
set -e

ENVIRONMENT=${1:-dev}

echo "Validating migrations for environment: $ENVIRONMENT"

DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-template_db}
DB_USER=${DB_USER:-template_user}
DB_PASSWORD=${DB_PASSWORD:-template_pass}
DB_SCHEMA=${DB_SCHEMA:-template_service}

echo "Using database: $DB_HOST:$DB_PORT/$DB_NAME schema=$DB_SCHEMA"

docker run --rm \
  -v "$(pwd):/liquibase/changelog" \
  -e "LIQUIBASE_COMMAND_USERNAME=$DB_USER" \
  -e "LIQUIBASE_COMMAND_PASSWORD=$DB_PASSWORD" \
  -e "LIQUIBASE_COMMAND_URL=jdbc:postgresql://$DB_HOST:$DB_PORT/$DB_NAME" \
  liquibase/liquibase:alpine \
