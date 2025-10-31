package repo

import (
	"context"
	"github.com/AlexAnd012/BookFinder/internal/data"
)

type BookRepo struct{ DB *Postgres }

func NewBookRepo(db *Postgres) *BookRepo { return &BookRepo{DB: db} }

// Create book
func (r *BookRepo) Create(ctx context.Context, b data.Book) (int64, error) {
	row := r.DB.Pool.QueryRow(ctx,
		`INSERT INTO books(title, language, pub_year, isbn) VALUES ($1,$2,$3,$4) RETURNING id`,
		b.Title, b.Language, b.PubYear, b.ISBN,
	)
	var id int64
	return id, row.Scan(&id)
}

// Get by ID
func (r *BookRepo) Get(ctx context.Context, id int64) (data.BookWithMeta, error) {
	var out data.BookWithMeta
	err := r.DB.Pool.QueryRow(ctx, `
    SELECT b.id, b.title, b.language, b.pub_year, b.isbn,
           
    --Для каждой книги b берём всех связанных авторов типа text[]
    --первый аргумент не NULL
    COALESCE((
      	SELECT array_agg(DISTINCT a.name) 
      	FROM book_author ba 
      		JOIN authors a ON a.id=ba.author_id 
      	WHERE ba.book_id=b.id), '{}') AS authors,
        
    -- собираем массив жанров для книги
    COALESCE((
    	SELECT array_agg(DISTINCT g.name) 
   	 	FROM book_genre bg 
   	 	    JOIN genres g ON g.id=bg.genre_id 
   	 	WHERE bg.book_id=b.id), '{}') AS genres,
        
    -- Средняя оценка книги по таблице reviews
      	(SELECT AVG(rating)::float8 
       	FROM reviews 
       	WHERE book_id=b.id) AS avg_rating
    
    FROM books b WHERE b.id=$1`, id).Scan(
		&out.ID, &out.Title, &out.Language, &out.PubYear, &out.ISBN, &out.Authors, &out.Genres, &out.AvgRating,
	)
	return out, err
}

// Search
func (r *BookRepo) Search(ctx context.Context, q *string, genre *string, yearFrom, yearTo *int, limit, offset int32) ([]data.BookWithMeta, error) {
	rows, err := r.DB.Pool.Query(ctx, `
    SELECT b.id, b.title, b.language, b.pub_year, b.isbn,
           
    --Для каждой книги b берём всех связанных авторов
    --первый аргумент не NULL
    COALESCE((
      	SELECT array_agg(DISTINCT a.name) 
      	FROM book_author ba 
      		JOIN authors a ON a.id=ba.author_id 
      	WHERE ba.book_id=b.id), '{}') AS authors,
        
	-- собираем массив жанров для книги
    COALESCE((
      	SELECT array_agg(DISTINCT g.name) 
      	FROM book_genre bg 
      	JOIN genres g ON g.id=bg.genre_id 
      	WHERE bg.book_id=b.id), '{}') AS genres,
        
    -- Средняя оценка книги по таблице reviews   
      (SELECT AVG(rating)::float8 
       FROM reviews 
       WHERE book_id=b.id) AS avg_rating
    
    FROM books b
    --$1 это строка поиска q
    --Если задано ILIKE (регистронезависимый) по подстроке идет склеивание
    WHERE ($1 IS NULL OR b.title ILIKE '%'||$1||'%')
    --$2 это жанр
    --Если жанр не задан то блок пропускаем. Иначе оставляем книги, для которых существует связь с этим жанром
      AND ($2 IS NULL OR EXISTS (
        SELECT 1 
        FROM book_genre bg 
            JOIN genres g ON g.id=bg.genre_id 
        WHERE bg.book_id=b.id AND g.name=$2
      ))
      --$3 / $4 это границы года издания
      AND ($3 IS NULL OR b.pub_year >= $3)
      AND ($4 IS NULL OR b.pub_year <= $4)
    ORDER BY avg_rating DESC NULLS LAST, 
        b.created_at DESC
    LIMIT $5 OFFSET $6`, q, genre, yearFrom, yearTo, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]data.BookWithMeta, 0)
	for rows.Next() {
		var b data.BookWithMeta
		if err := rows.Scan(&b.ID, &b.Title, &b.Language, &b.PubYear, &b.ISBN, &b.Authors, &b.Genres, &b.AvgRating); err != nil {
			return nil, err
		}
		res = append(res, b)
	}
	return res, rows.Err()
}
