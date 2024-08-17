package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/eclipsemode/go-bot-tg-helper/storage"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// New creates new SQLite storage
func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("cannot open db: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("cannot ping db: %v", err)
	}

	return &Storage{db: db}, nil
}

// Save saves page to storage.
func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := `INSERT INTO pages (url, user_name) VALUES (?, ?)`

	_, err := s.db.ExecContext(ctx, q, p.URL, p.UserName)
	if err != nil {
		return fmt.Errorf("cannot execute query: %v", err)
	}

	return nil
}

// PickRandom pick random page from storage.
func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, q, userName).Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cannot execute query: %w", err)
	}

	return &storage.Page{URL: url}, nil
}

// Remove removes page from storage.
func (s *Storage) Remove(ctx context.Context, page *storage.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND user_name = ?`

	_, err := s.db.ExecContext(ctx, q, page.URL, page.UserName)
	if err != nil {
		return fmt.Errorf("cannot execute query: %w", err)
	}
	return nil
}

// IsExists checks if page exists in storage.
func (s *Storage) IsExists(ctx context.Context, page *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, page.URL, page.UserName).Scan(&count); err != nil {
		return false, fmt.Errorf("cannot execute query: %w", err)
	}

	return count > 0, nil
}

// Init creates row pages if not exists.
func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`
	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("cannot execute query: %w", err)
	}

	return nil
}
