#!/usr/bin/env bash
set -euo pipefail
: "${DB_USER:=postgres}" "${DB_PASSWORD:=postgres}" "${DB_NAME:=book_finder}"
psql postgresql://$DB_USER:$DB_PASSWORD@localhost:5432/$DB_NAME?sslmode=disable -f migrations/0002_sample_data.sql
