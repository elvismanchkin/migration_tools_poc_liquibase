#!/bin/sh
set -e

echo "Waiting for database to be ready..."
until pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME; do
  echo "Database not ready yet, waiting..."
  sleep 2
done

echo "Database is ready!"
echo "Starting template service..."
exec /app/template-service