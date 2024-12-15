package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	// In this place rise the error
	const op = "database.postgres.New"

	db, err := sql.Open("postgres", storagePath)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "database.postgres.SaveUrl"
	println(urlToSave)
	println(alias)

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES ($1, $2)")

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// TODO: Пофиксить это
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil

}

func (s *Storage) GetUrl(alias string) (string, error) {
	const op = "database.postgres.GetUrl"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = $1")

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	var resultUrl string
	err = stmt.QueryRow(alias).Scan(&resultUrl)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resultUrl, nil

}

func (s *Storage) DeleteUrl(alias string) error {
	const op = "database.postgres.DeleteUrl"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = $1")

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(alias)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
