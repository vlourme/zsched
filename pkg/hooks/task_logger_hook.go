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
			status SYMBOL,
			task_name VARCHAR,
			parent_id UUID,
			state VARCHAR,
			iterations INTEGER,
			published_at TIMESTAMP,
			started_at TIMESTAMP,
			ended_at TIMESTAMP,
			last_error VARCHAR
		) TIMESTAMP(published_at) 
		PARTITION BY DAY TTL 7 DAYS
		DEDUP UPSERT KEYS (published_at, task_id)
	`)
	if err != nil {
		log.Fatalf("failed to create task logs table: %v", err)
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
	`, s.TaskID, s.Status, task.Name(), s.ParentID, parameters, s.Iterations, s.InitializedAt, s.StartedAt, time.Now(), s.LastError)
	if err != nil {
		log.Printf("failed to update task in storage: %v", err)
		return err
	}

	return nil
}
