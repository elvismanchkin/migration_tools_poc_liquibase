#!/bin/bash
set -e

# wrapper for the liquibase cli
# usage: ./liquibase-api.sh <command> [environment]
#
# examples:
#   ./liquibase-api.sh update dev
#   ./liquibase-api.sh rollback prod 1
#   ./liquibase-api.sh status

COMMAND=${1:-update}
ENVIRONMENT=${2:-dev}
EXTRA_ARGS=${@:3}

DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-template_db}
DB_USER=${DB_USER:-template_user}
DB_PASSWORD=${DB_PASSWORD:-template_pass}
DB_SCHEMA=${DB_SCHEMA:-template_service}

echo "Executing Liquibase command: $COMMAND"
echo "Environment: $ENVIRONMENT"
echo "Using database: $DB_HOST:$DB_PORT/$DB_NAME schema=$DB_SCHEMA"

ADDITIONAL_ARGS=""
if [ "$COMMAND" = "rollback" ]; then
  if [ -z "$EXTRA_ARGS" ]; then
    echo "Error: Rollback requires a tag, count, or date parameter"
    exit 1
  fi
  ADDITIONAL_ARGS="$EXTRA_ARGS"
elif [ "$COMMAND" = "update" ]; then
  ADDITIONAL_ARGS="--contexts=$ENVIRONMENT"
fi

docker run --rm \
  -v "$(pwd):/liquibase/changelog" \
  -e "LIQUIBASE_COMMAND_USERNAME=$DB_USER" \
  -e "LIQUIBASE_COMMAND_PASSWORD=$DB_PASSWORD" \
  -e "LIQUIBASE_COMMAND_URL=jdbc:postgresql://$DB_HOST:$DB_PORT/$DB_NAME" \
  -e "LIQUIBASE_DEFAULTS_FILE=changelog/liquibase.properties" \
  liquibase/liquibase:alpine \
