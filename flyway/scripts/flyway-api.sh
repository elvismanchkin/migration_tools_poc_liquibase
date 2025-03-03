#!/bin/bash
set -e

# wrapper for the flyway cli
# usage: ./flyway-api.sh <command> [environment]
#
# Examples:
#   ./flyway-api.sh migrate dev
#   ./flyway-api.sh repair
#   ./flyway-api.sh info

COMMAND=${1:-migrate}
ENVIRONMENT=${2:-dev}
EXTRA_ARGS=${@:3}

DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-template_db}
DB_USER=${DB_USER:-template_user}
DB_PASSWORD=${DB_PASSWORD:-template_pass}
DB_SCHEMA=${DB_SCHEMA:-template_service}

echo "Executing Flyway command: $COMMAND"
echo "Environment: $ENVIRONMENT"
echo "Using database: $DB_HOST:$DB_PORT/$DB_NAME schema=$DB_SCHEMA"

ENV_PLACEHOLDER=""
[ -n "$ENVIRONMENT" ] && ENV_PLACEHOLDER="-placeholders.environment=$ENVIRONMENT"

docker run --rm \
  -v "$(pwd)/sql:/flyway/sql" \
  -e "FLYWAY_URL=jdbc:postgresql://$DB_HOST:$DB_PORT/$DB_NAME" \
  -e "FLYWAY_USER=$DB_USER" \
  -e "FLYWAY_PASSWORD=$DB_PASSWORD" \
  -e "FLYWAY_DEFAULT_SCHEMA=$DB_SCHEMA" \
