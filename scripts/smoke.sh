#!/usr/bin/env bash

# Проверяем работает ли сервис

#-e — выйти при любой ошибке команды
#-u — ошибка, если обращаемся к неопределённой переменной
#-o pipefail — если команда в пайпе упала, весь пайп считается ошибкой
set -euo pipefail

# запрашиваем /health
curl -s localhost:8080/health | jq . || true
# запрос проверки книг
curl -s "localhost:8080/v1/books?q=Hobbit&limit=10" | jq . || true
