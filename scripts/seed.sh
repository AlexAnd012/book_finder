#!/usr/bin/env bash

#-e — выйти при любой ошибке команды
#-u — ошибка, если обращаемся к неопределённой переменной
#-o pipefail — если команда в пайпе упала, весь пайп считается ошибкой
set -euo pipefail

# читаем из окружения
: "${DB_USER:=postgres}" "${DB_PASSWORD:=postgres}" "${DB_NAME:=book_finder}"

# Запускаем psql и выполняем сидирование данными
psql postgresql://$DB_USER:$DB_PASSWORD@localhost:5432/$DB_NAME?sslmode=disable -f migrations/0002_sample_data.sql
