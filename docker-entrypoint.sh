#!/bin/sh
set -e

echo "Running database migrations..."
/app/prisoner migrate up

echo "Starting server..."
exec /app/prisoner server
