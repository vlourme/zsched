package zsched

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/vlourme/zsched/pkg/storage"
)

type taskLogger[T any] struct {
	storage storage.Storage
	pending chan pendingTask
}

type pendingTask struct {
	TaskID        uuid.UUID
	Status        stateStatus
	TaskName      string
	ParentID      uuid.UUID
	Parameters    string
	Iterations    int
	InitializedAt time.Time
	StartedAt     time.Time
	EndedAt       time.Time
	LastError     string
}

func NewTaskLogger[T any](storage storage.Storage) (*taskLogger[T], error) {
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

	tl := &taskLogger[T]{
		storage: storage,
		pending: make(chan pendingTask, 1000),
	}
	go tl.worker()
	return tl, nil
}

func (h *taskLogger[T]) worker() {
	batch := h.storage.NewBatch()
	interval := time.NewTicker(time.Second)
	shouldFlush := make(chan bool, 1)

	for {
		select {
		case pending := <-h.pending:
			batch.Add(
				`
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
				`,
				pending.TaskID,
				pending.Status,
				pending.TaskName,
				pending.ParentID,
				pending.Parameters,
				pending.Iterations,
				pending.InitializedAt,
				pending.StartedAt,
				pending.EndedAt,
				pending.LastError,
			)

			if batch.Size() >= 5000 {
				select {
				case shouldFlush <- true:
				default:
				}
			}
		case <-interval.C:
			select {
			case shouldFlush <- true:
			default:
			}
		case <-shouldFlush:
			if err := batch.Execute(); err != nil {
				log.Printf("failed to flush tasks to storage: %v", err)
			}
			batch = h.storage.NewBatch()
		}
	}
}

func (h *taskLogger[T]) LogTasks(task *Task[T], state *State) error {
	parameters, err := state.EncodeParameters()
	if err != nil {
		log.Printf("failed to encode parameters: %v", err)
		return err
	}

	pending := pendingTask{
		TaskID:        state.TaskID,
		Status:        state.Status,
		TaskName:      task.Name(),
		ParentID:      state.ParentID,
		Parameters:    parameters,
		Iterations:    state.Iterations,
		InitializedAt: state.InitializedAt,
		StartedAt:     state.StartedAt,
	}

	if state.Status == StatusSuccess || state.Status == StatusFailed {
		pending.EndedAt = time.Now()
	}

	h.pending <- pending

	return nil
}
