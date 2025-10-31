INSERT INTO authors(name) VALUES ('J. R. R. Tolkien') ON CONFLICT DO NOTHING;
INSERT INTO authors(name) VALUES ('Ursula K. Le Guin') ON CONFLICT DO NOTHING;

INSERT INTO genres(name) VALUES ('Fantasy') ON CONFLICT DO NOTHING;
INSERT INTO genres(name) VALUES ('Sci-Fi')  ON CONFLICT DO NOTHING;

INSERT INTO books(title, language, pub_year, isbn)
VALUES ('The Hobbit', 'en', 1937, '9780547928227')
ON CONFLICT (isbn) DO NOTHING;

INSERT INTO book_author(book_id, author_id)
SELECT b.id, a.id
FROM books b, authors a
WHERE b.title='The Hobbit' AND a.name='J. R. Tolkien'
ON CONFLICT DO NOTHING;

INSERT INTO book_genre(book_id, genre_id)
SELECT b.id, g.id
FROM books b, genres g
WHERE b.title='The Hobbit' AND g.name='Fantasy'
ON CONFLICT DO NOTHING;
