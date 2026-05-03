#!/bin/bash

# Database setup script for NowYouSeeMe Event Sourcing

set -e

DB_NAME="${DB_NAME:-nowyouseeme}"
DB_USER="${DB_USER:-postgres}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"

echo "Setting up database: $DB_NAME"

# Check if PostgreSQL is running
if ! pg_isready -h "$DB_HOST" -p "$DB_PORT" > /dev/null 2>&1; then
  echo "Error: PostgreSQL is not running on $DB_HOST:$DB_PORT"
  exit 1
fi

# Drop existing database if exists
echo "Dropping existing database (if exists)..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -c "DROP DATABASE IF EXISTS $DB_NAME;" postgres

# Create fresh database
echo "Creating database..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -c "CREATE DATABASE $DB_NAME;" postgres

# Run migration
echo "Running migration..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f backend/migrations/001_create_event_sourcing_schema.sql

echo "Database setup complete!"
