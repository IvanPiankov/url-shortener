package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/storage"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3" // init sqlite3 driver
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	// In this place rise the error
	const op = "database.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(ctx context.Context, urlToSave string, alias string) (int64, error) {
	const op = "database.sqlite.SaveUrl"

	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO url(url, alias) VALUES (?, ?)")

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, urlToSave, alias)
	if err != nil {
		// TODO: refactor this
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLAlreadyExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil

}

func (s *Storage) GetUrl(ctx context.Context, alias string) (string, error) {
	const op = "database.sqlite.GetUrl"

	stmt, err := s.db.PrepareContext(ctx, "SELECT url FROM url WHERE alias = ?")

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	var resultUrl string
	err = stmt.QueryRowContext(ctx, alias).Scan(&resultUrl)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resultUrl, nil

}

func (s *Storage) DeleteUrl(ctx context.Context, alias string) error {
	const op = "database.sqlite.DeleteUrl"

	stmt, err := s.db.PrepareContext(ctx, "DELETE FROM url WHERE alias = ?")

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, alias)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
