package logger

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vlourme/zsched/pkg/storage"
)

type Logger = *logrus.Entry

func NewLogger(storage storage.Storage) Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.Kitchen,
		FullTimestamp:   true,
	})

	if storage.Name() == "questdb" {
		logger.AddHook(NewQuestDBHook(storage))
	}

	return logrus.NewEntry(logger)
}
