#!/bin/sh
set -e

DB_PATH="${DATABASE_PATH:-data/trama.db}"

echo "Recreating database from SQL files..."
rm -f "$DB_PATH"

for f in /app/localdb/*.sql; do
    echo "  Running $(basename "$f")..."
    sqlite3 "$DB_PATH" < "$f"
done

echo "Database ready."
exec ./trama
