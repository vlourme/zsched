package storage

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type TimescaleDBStorage struct {
	conn *sql.DB
}

// NewTimescaleDBStorage creates a new TimescaleDB storage.
func NewTimescaleDBStorage(dsn string) (*TimescaleDBStorage, error) {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &TimescaleDBStorage{conn}, nil
}

func (s *TimescaleDBStorage) Close() error {
	return s.conn.Close()
}

func (s *TimescaleDBStorage) Exec(query string, args ...any) (sql.Result, error) {
	return s.conn.Exec(query, args...)
}

func (s *TimescaleDBStorage) Query(query string, args ...any) (*sql.Rows, error) {
	return s.conn.Query(query, args...)
}

func (s *TimescaleDBStorage) QueryRow(query string, args ...any) *sql.Row {
	return s.conn.QueryRow(query, args...)
}

func (s *TimescaleDBStorage) Connection() (*sql.Conn, error) {
	return s.conn.Conn(context.Background())
}

func (s *TimescaleDBStorage) Name() string {
	return "timescaledb"
}
