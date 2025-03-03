# Liquibase Migration Scripts

This directory contains utility scripts to help manage database migrations using Liquibase.

## Available Scripts

### liquibase-migrate.sh

Applies all pending migrations to the database.

```bash
./liquibase-migrate.sh [environment]
```

- `environment`: The environment context to use (default: `dev`)

Example:
```bash
./liquibase-migrate.sh prod
```

### liquibase-validate.sh

Validates the changelog without making any changes to the database.

```bash
./liquibase-validate.sh [environment]
```

Example:
```bash
./liquibase-validate.sh dev
```

### liquibase-api.sh

A general-purpose wrapper for Liquibase commands.

```bash
./liquibase-api.sh <command> [environment] [extra-args]
```

Examples:
```bash
./liquibase-api.sh update dev        # Apply migrations for dev environment
