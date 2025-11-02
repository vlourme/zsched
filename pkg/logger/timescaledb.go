package logger

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/vlourme/zsched/pkg/storage"
)

type TimescaleDBHook struct {
	storage storage.Storage
}

func NewTimescaleDBHook(storage storage.Storage) *TimescaleDBHook {
	_, err := storage.Exec(`
		CREATE TABLE IF NOT EXISTS logs (
			task_id UUID,
			state_id UUID,
			level VARCHAR(10),
			message TEXT,
			data JSONB,
			logged_at TIMESTAMPTZ,
			PRIMARY KEY (task_id, state_id, logged_at)
		)
		WITH (
			tsdb.hypertable,
			tsdb.partition_column='logged_at',
			tsdb.orderby='logged_at DESC'
		)
	`)
	if err != nil {
		log.Fatalf("failed to create logs table: %v", err)
	}

	_, err = storage.Exec(
		`SELECT add_retention_policy('logs', drop_after => INTERVAL '7 days', if_not_exists => true)`,
	)
	if err != nil {
		log.Fatalf("failed to create retention policy: %v", err)
	}

	return &TimescaleDBHook{storage}
}

func (h *TimescaleDBHook) Fire(entry *logrus.Entry) error {
	taskId, ok := entry.Data["task_id"]
	if !ok {
		return nil
	}
	delete(entry.Data, "task_id")

	stateId, ok := entry.Data["state_id"]
	if !ok {
		return nil
	}
	delete(entry.Data, "state_id")

	data, err := json.Marshal(entry.Data)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO logs (task_id, state_id, level, message, data, logged_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = h.storage.Exec(query, taskId, stateId, entry.Level, entry.Message, data, entry.Time)
	if err != nil {
		return errors.Join(errors.New("failed to insert log"), err)
	}

	return err
}

func (h *TimescaleDBHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
