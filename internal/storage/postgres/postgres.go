package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"url-shortener/internal/storage"
)

type Storage struct {
	db *pgx.Conn
}

func New(storagePath string) (*Storage, error) {
	// In this place rise the error
	const op = "database.postgres.New"

	db, err := pgx.Connect(context.Background(), storagePath)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(ctx context.Context, urlToSave string, alias string) (int64, error) {
	const op = "database.postgres.SaveUrl"

	query := `INSERT INTO url(url, alias) VALUES (@url, @alias) RETURNING id;`
	args := pgx.NamedArgs{
		"url":   urlToSave,
		"alias": alias,
	}

	var aliasId int64
	err := s.db.QueryRow(ctx, query, args).Scan(&aliasId)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return aliasId, nil

}

func (s *Storage) GetUrl(ctx context.Context, alias string) (string, error) {
	const op = "database.postgres.GetUrl"

	query := `SELECT url FROM url WHERE alias = $1`

	var resultUrl string

	err := s.db.QueryRow(ctx, query, alias).Scan(&resultUrl)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}

	return resultUrl, nil

}

func (s *Storage) DeleteUrl(ctx context.Context, alias string) error {
	const op = "database.postgres.DeleteUrl"

	query := `DELETE FROM url WHERE alias = $1`

	_, err := s.db.Exec(ctx, query, alias)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
