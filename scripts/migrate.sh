#!/usr/bin/env bash
set -euo pipefail
: "${DB_USER:=postgres}" "${DB_PASSWORD:=postgres}" "${DB_NAME:=book_finder}"
PSQL="psql postgresql://$DB_USER:$DB_PASSWORD@localhost:5432/$DB_NAME?sslmode=disable"

# создать БД, если нет
createdb -h localhost -U "$DB_USER" "$DB_NAME" 2>/dev/null || true

for f in migrations/*.sql; do
  echo "Applying $f";
  $PSQL -v ON_ERROR_STOP=1 -f "$f";
done
