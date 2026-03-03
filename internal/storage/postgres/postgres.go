package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/storage"

	"github.com/jackc/pgx/v5/pgconn"
)

// Реализация хранилища PostgreSQL
type Storage struct {
	db *sql.DB
}

// Конструктор для Storage
func New(dbPath string) (*Storage, error) {
	const op = "storage.postgres.New"

	// Подключение к бд
	db, err := sql.Open("pgx", dbPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Проверка подключения к бд
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Создать табл, если еще нет
    stmt, err := db.Prepare(`
    CREATE TABLE IF NOT EXISTS url(
        id INTEGER PRIMARY KEY,
        alias TEXT NOT NULL UNIQUE,
        url TEXT NOT NULL);
    CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
    `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// Сохранение URL и алиаса
func (st *Storage) SaveURL(urlToSave string, alias string) (error) {
	const op = "storage.postgres.SaveURL"

	// Подготовка запроса
	stmt, err := st.db.Prepare(`INSERT INTO url (url, alias) VALUES ($1, $2) RETURNING id`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	
	// Выполнение запроса
	var id int64
	err = stmt.QueryRow(urlToSave, alias).Scan(&id)	
	if err != nil {
		// Проверка алиаса на уникальность
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Получение URL по алиасу
func (st *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"

	// Подготовка запроса
	stmt, err := st.db.Prepare(`SELECT url FROM url WHERE alias = $1`)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Выполнение запроса
	var url string
	err = stmt.QueryRow(alias).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}
