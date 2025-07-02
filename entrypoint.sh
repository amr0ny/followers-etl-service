#!/bin/sh

set -e

echo "Running database migrations..."
migrate -path ./migrations -database "${DSN}" up || { echo "Migration failed"; exit 1; }

echo "Starting application..."
exec "$@"