package hooks

import (
	"log"
	"time"

	"github.com/vlourme/zsched"
	"github.com/vlourme/zsched/pkg/storage"
)

type TaskLoggerHook struct {
	storage storage.Storage
}

func (h *TaskLoggerHook) Initialize(storage storage.Storage) error {
	h.storage = storage

	_, err := storage.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			task_id UUID,
			status VARCHAR(10),
			task_name VARCHAR(128),
			parent_id UUID,
			state JSONB,
			iterations INTEGER,
			published_at TIMESTAMPTZ,
			started_at TIMESTAMPTZ,
			ended_at TIMESTAMPTZ,
			last_error TEXT,
			PRIMARY KEY (task_id, published_at)
		) WITH (
			tsdb.hypertable,
			tsdb.partition_column='published_at',
			tsdb.orderby='published_at DESC'
		)
	`)
	if err != nil {
		log.Fatalf("failed to create task logs table: %v", err)
	}

	_, err = storage.Exec(
		`SELECT add_retention_policy('tasks', drop_after => INTERVAL '7 days', if_not_exists => true)`,
	)
	if err != nil {
		log.Fatalf("failed to create retention policy for tasks table: %v", err)
	}

	return nil
}

func (h *TaskLoggerHook) BeforeExecute(task zsched.AnyTask, s *zsched.State) error {
	parameters, err := s.EncodeParameters()
	if err != nil {
		log.Printf("failed to encode parameters: %v", err)
		return err
	}

	_, err = h.storage.Exec(`
		INSERT INTO tasks (task_id, status, task_name, parent_id, state, iterations, published_at, started_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (task_id, published_at)
		DO UPDATE SET
			status = $2,
			task_name = $3,
			parent_id = $4,
			state = $5,
			iterations = $6,
			published_at = $7,
			started_at = $8
	`, s.TaskID, s.Status, task.Name(), s.ParentID, parameters, s.Iterations, s.InitializedAt, s.StartedAt)
	if err != nil {
		log.Printf("failed to insert task into storage: %v", err)
		return err
	}

	return nil
}

func (h *TaskLoggerHook) AfterExecute(task zsched.AnyTask, s *zsched.State) error {
	parameters, err := s.EncodeParameters()
	if err != nil {
		log.Printf("failed to encode parameters: %v", err)
		return err
	}

	_, err = h.storage.Exec(`
		INSERT INTO tasks (task_id, status, task_name, parent_id, state, iterations, published_at, started_at, ended_at, last_error)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (task_id, published_at)
		DO UPDATE SET
			status = $2,
			task_name = $3,
			parent_id = $4,
			state = $5,
			iterations = $6,
			published_at = $7,
			started_at = $8,
			ended_at = $9,
			last_error = $10
	`, s.TaskID, s.Status, task.Name(), s.ParentID, parameters, s.Iterations, s.InitializedAt, s.StartedAt, time.Now(), s.LastError)
	if err != nil {
		log.Printf("failed to update task in storage: %v", err)
		return err
	}

	return nil
}
