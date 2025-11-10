package storage

import "database/sql"

type Storage interface {
	// DB returns the database
	Connection() (*sql.Conn, error)

	// Exec executes a query
	Exec(query string, args ...any) (sql.Result, error)

	// Query executes a query and returns a rows
	Query(query string, args ...any) (*sql.Rows, error)

	// NewBatch creates a new batch
	NewBatch() Batch

	// QueryRow executes a query and returns a single row
	QueryRow(query string, args ...any) *sql.Row

	// Close closes the storage
	Close() error

	// Name returns the name of the storage
	Name() string
}

// Batch is a batch of queries to be executed
type Batch interface {
	// Add adds a query to the batch
	Add(query string, args ...any) error

	// Size returns the number of queries in the batch
	Size() int

	// Execute executes the batch
	Execute() error
}
