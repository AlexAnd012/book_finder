#!/usr/bin/env bash

#-e — выйти при любой ошибке команды
#-u — ошибка, если обращаемся к неопределённой переменной
#-o pipefail — если команда в пайпе упала, весь пайп считается ошибкой
set -euo pipefail

# читаем из окружения
: "${DB_USER:=postgres}" "${DB_PASSWORD:=postgres}" "${DB_NAME:=book_finder}"
# команда для psql с полным DSN
PSQL="psql postgresql://$DB_USER:$DB_PASSWORD@localhost:5432/$DB_NAME?sslmode=disable"

# создаем БД, если нет
createdb -h localhost -U "$DB_USER" "$DB_NAME" 2>/dev/null || true

# цикл по всем SQL-файлам в папке migrations/ и выполняем SQL из файла
for f in migrations/*.sql; do
  echo "Applying $f";
  $PSQL -v ON_ERROR_STOP=1 -f "$f";
done
