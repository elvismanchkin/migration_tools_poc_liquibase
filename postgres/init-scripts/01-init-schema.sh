#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Create schemas
    CREATE SCHEMA IF NOT EXISTS template_service;
    CREATE SCHEMA IF NOT EXISTS audit;

    -- Set search path
    ALTER DATABASE $POSTGRES_DB SET search_path TO template_service, public;

    -- Create extension for UUIDs
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

    -- Grant permissions
    GRANT ALL PRIVILEGES ON SCHEMA template_service TO $POSTGRES_USER;
    GRANT ALL PRIVILEGES ON SCHEMA audit TO $POSTGRES_USER;
EOSQL

echo "PostgreSQL initialized with template_service and audit schemas"