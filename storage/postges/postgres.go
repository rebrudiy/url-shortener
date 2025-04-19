package postges

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(connString string) (*Storage, error) {
	const op = "storage.postgres.New"
	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = db.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS url(
				id SERIAL PRIMARY KEY,
				alias TEXT NOT NULL UNIQUE,
				url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
		`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil

}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.postgres.SaveURL"
	query, args, err := sq.Insert("url").
		Columns("alias", "url").
		Values(alias, urlToSave).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	var id int64
	err = s.db.QueryRow(context.Background(), query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s : %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GerURL"
	query, args, err := sq.
		Select("url").
		From("url").
		Where(sq.Eq{"alias": alias}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("%s: failed to build query : %w", op, err)
	}

	var url string
	err = s.db.QueryRow(context.Background(), query, args...).Scan(&url)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, err
}

//TODO: implement delete url
