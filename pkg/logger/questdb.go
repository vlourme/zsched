package logger

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/vlourme/zsched/pkg/storage"
)

type QuestDBHook struct {
	storage storage.Storage
}

func NewQuestDBHook(storage storage.Storage) *QuestDBHook {
	_, err := storage.Exec(`
		CREATE TABLE IF NOT EXISTS logs (
			task_id UUID,
			state_id UUID,
			level SYMBOL,
			message VARCHAR,
			data VARCHAR,
			logged_at TIMESTAMP
		)
		TIMESTAMP(logged_at)
		PARTITION BY DAY TTL 7 DAYS
	`)
	if err != nil {
		log.Fatalf("failed to create logs table: %v", err)
	}

	return &QuestDBHook{storage}
}

func (h *QuestDBHook) Fire(entry *logrus.Entry) error {
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

func (h *QuestDBHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
