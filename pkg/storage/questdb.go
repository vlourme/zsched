package storage

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type QuestDBStorage struct {
	conn *sql.DB
}

// NewQuestDBStorage creates a new QuestDB storage
// QuestDB is based on PostgreSQL protocol.
func NewQuestDBStorage(dsn string) (*QuestDBStorage, error) {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &QuestDBStorage{conn}, nil
}

func (s *QuestDBStorage) Close() error {
	return s.conn.Close()
}

func (s *QuestDBStorage) Exec(query string, args ...any) (sql.Result, error) {
	return s.conn.Exec(query, args...)
}

func (s *QuestDBStorage) Query(query string, args ...any) (*sql.Rows, error) {
	return s.conn.Query(query, args...)
}

func (s *QuestDBStorage) QueryRow(query string, args ...any) *sql.Row {
	return s.conn.QueryRow(query, args...)
}

func (s *QuestDBStorage) Connection() (*sql.Conn, error) {
	return s.conn.Conn(context.Background())
}

func (s *QuestDBStorage) Name() string {
	return "questdb"
}
