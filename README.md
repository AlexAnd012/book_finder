# Book Finder API (Backend)
## 1) Цель

Создать сервис каталогизации и поиска книг по нескольким критериям с импортом данных, отзывами и рейтингами. Сервис предоставляет REST API, работает с PostgreSQL, поддерживает фоновые задачи через каналы и корректно обрабатывает контекст и отмену.

## 2) Технологический стек 

Язык: Go 1.24

Веб-фреймворк: chi 

БД: PostgreSQL 17

Контейнеризация: Docker + Docker Compose

Миграции: golang-migrate (CLI) или встроенные миграции

Тесты: Go test, покрытия не ниже 75% по package internal/...

Логи: structured JSON 

## 3) Область функционала (MVP)

Сущности:

Пользователь (регистрация/логин, роли: user, admin)

Автор

Жанр

Книга

Отзыв/Рейтинг (1–5)

CRUD

Книги, Авторы, Жанры 

Отзывы — create/list для авторизованных пользователей

Поиск

Поиск книг по: q (title/author, ILIKE + триграммы), author, genre, language, year_from, year_to

Сортировка: -rating, rating, -created, created, title

Пагинация: limit/offset 

Аутентификация/Авторизация

JWT (HS256), обновление через refresh-токен

## 4) REST API 
### Книги

POST /v1/books (admin)
body: {title, language?, pub_year?, isbn?, authors: [name|id], genres:[name|id]}

GET /v1/books
query:
q?, author?, genre?, language?, year_from?, year_to?, sort?=-rating|rating|created|-created|title, limit?<=100, offset?
resp: {items:[...BookDTO], total, limit, offset}

GET /v1/books/{id}
resp: BookDTO с authors[], genres[], avg_rating, reviews_count

PATCH /v1/books/{id} (admin)

DELETE /v1/books/{id} (admin)

### Авторы

POST /v1/authors (admin)

GET /v1/authors?name=&limit=&offset=

GET /v1/authors/{id}

PATCH /v1/authors/{id} (admin)

DELETE /v1/authors/{id} (admin)

### Жанры

POST /v1/genres (admin)

GET /v1/genres?name=&limit=&offset=

DELETE /v1/genres/{id} (admin)

### Отзывы

POST /v1/reviews (user)
body: {book_id, rating (1..5), text?}

GET /v1/books/{id}/reviews?limit=&offset=

### Импорт

POST /v1/import (admin)
form-data: file (CSV) или JSON: {url:"https://..."}
resp: {job_id}

GET /v1/import/{job_id} → {status: queued|running|done|failed, stats:{inserted, updated, skipped, errors}}


## 5) Структура репозитория
book-finder/
├─ cmd/main.go  
├─ internal/config/config.go  
├─ internal/httpserver/router.go 
├─ internal/httpserver/router_test.go  
├─ internal/httpserver/middleware_log.go  
├─ internal/handlers/books.go  
├─ internal/handlers/health.go  
├─ internal/handlers/books_test.go  
├─ internal/handlers/health_test.go  
├─ internal/repo/postgres.go  
├─ internal/repo/book_repo.go  
├─ internal/logging/logger.go 
├─ internal/data/data.go 
├─ migrations/  
│ ├─ 0001_init.sql  
│ └─ 0002_sample_data.sql  
├─ scripts/  
│ ├─ migrate.sh  
│ ├─ seed.sh  
│ └─ smoke.sh  
├─ docker-compose.yml  
├─ .env
└─ go.mod  

## 6) Скрипты Bash 

scripts/migrate.sh — применить миграции 

scripts/seed.sh — базовые данные 

scripts/smoke.sh — curl healthz, create user
