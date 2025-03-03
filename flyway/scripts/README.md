# Flyway Migration Scripts

This directory contains utility scripts to help manage database migrations using Flyway.

## Available Scripts

### flyway-migrate.sh

Applies all pending migrations to the database.

```bash
./flyway-migrate.sh [environment]
```

- `environment`: The environment to use for placeholders (default: `dev`)

Example:

```bash
./flyway-migrate.sh prod
```

### flyway-validate.sh

Validates the migrations without making any changes to the database.

```bash
./flyway-validate.sh [environment]
```

Example:

```bash
./flyway-validate.sh dev
```

### flyway-api.sh

A general-purpose wrapper for Flyway commands.

```bash
./flyway-api.sh <command> [environment] [extra-args]
```

Examples:

```bash
./flyway-api.sh migrate dev       # Apply migrations for dev environment
