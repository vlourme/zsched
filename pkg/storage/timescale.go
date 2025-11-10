package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// TimescaleDBStorage is the storage for the TimescaleDB database
type TimescaleDBStorage struct {
	pool *pgxpool.Pool
	conn *sql.DB
}

// TimescaleDBBatch is the batch for the TimescaleDB database
type TimescaleDBBatch struct {
	conn  *pgxpool.Pool
	batch *pgx.Batch
}

// NewTimescaleDBStorage creates a new TimescaleDB storage.
func NewTimescaleDBStorage(dsn string) (*TimescaleDBStorage, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &TimescaleDBStorage{
		pool: pool,
		conn: stdlib.OpenDBFromPool(pool),
	}, nil
}

func (s *TimescaleDBStorage) Close() error {
	return s.conn.Close()
}

func (s *TimescaleDBStorage) NewBatch() Batch {
	return &TimescaleDBBatch{
		conn:  s.pool,
		batch: &pgx.Batch{},
	}
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

func (b *TimescaleDBBatch) Add(query string, args ...any) error {
	b.batch.Queue(query, args...)
	return nil
}

func (b *TimescaleDBBatch) Size() int {
	return b.batch.Len()
}

func (b *TimescaleDBBatch) Execute() error {
	if e := b.conn.SendBatch(context.Background(), b.batch).Close(); e != nil {
		return errors.Join(errors.New("failed to execute batch"), e)
	}

	return nil
}
