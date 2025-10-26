package hooks

import (
	"log"
	"time"

	"github.com/vlourme/zsched/pkg/state"
	"github.com/vlourme/zsched/pkg/storage"
	"github.com/vlourme/zsched/pkg/task"
)

type TaskLoggerHook struct {
	storage storage.Storage
}

func NewTaskLoggerHook(storage storage.Storage) *TaskLoggerHook {
	_, err := storage.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			task_id UUID,
			task_name VARCHAR,
			parent_id UUID,
			state VARCHAR,
			iterations INTEGER,
			started_at TIMESTAMP,
			ended_at TIMESTAMP,
			last_error VARCHAR
		) TIMESTAMP(started_at) 
		PARTITION BY DAY TTL 7 DAYS
		DEDUP UPSERT KEYS (started_at, task_id)
	`)
	if err != nil {
		log.Fatalf("failed to create task logs table: %v", err)
	}

	return &TaskLoggerHook{
		storage: storage,
	}
}

func (h *TaskLoggerHook) BeforeExecute(task *task.Task, s *state.State) error {
	if s.Status != state.StatusPending {
		return nil
	}

	parameters, err := s.EncodeParameters()
	if err != nil {
		log.Printf("failed to encode parameters: %v", err)
		return err
	}

	_, err = h.storage.Exec(`
		INSERT INTO tasks (task_id, task_name, parent_id, state, iterations, started_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, s.TaskID, task.Name(), s.ParentID, parameters, s.Iterations, s.InitializedAt)
	if err != nil {
		log.Printf("failed to insert task into storage: %v", err)
		return err
	}

	return nil
}

func (h *TaskLoggerHook) AfterExecute(task *task.Task, s *state.State) error {
	if s.Status != state.StatusSuccess && s.Status != state.StatusFailed {
		return nil
	}

	parameters, err := s.EncodeParameters()
	if err != nil {
		log.Printf("failed to encode parameters: %v", err)
		return err
	}

	_, err = h.storage.Exec(`
		INSERT INTO tasks (task_id, task_name, parent_id, state, iterations, started_at, ended_at, last_error)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, s.TaskID, task.Name(), s.ParentID, parameters, s.Iterations, s.InitializedAt, time.Now(), s.LastError)
	if err != nil {
		log.Printf("failed to update task in storage: %v", err)
		return err
	}

	return nil
}
