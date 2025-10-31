#!/usr/bin/env bash
set -euo pipefail
curl -s localhost:8080/health | jq . || true
curl -s "localhost:8080/v1/books?q=Hobbit&limit=10" | jq . || true
