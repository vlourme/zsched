package hooks

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vlourme/zsched"
	"github.com/vlourme/zsched/pkg/storage"
)

type PrometheusHook struct {
	taskCounter       *prometheus.CounterVec
	durationHistogram *prometheus.HistogramVec
}

func (h *PrometheusHook) Initialize(storage storage.Storage) error {
	h.taskCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "scheduler_tasks_total",
		Help: "Total number of tasks executed",
	}, []string{"task_name", "status"})

	h.durationHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "scheduler_task_duration_seconds",
		Help: "Duration of tasks in seconds",
	}, []string{"task_name", "status"})

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	return nil
}

func (h *PrometheusHook) BeforeExecute(task zsched.AnyTask, s *zsched.State) error {
	return nil
}

func (h *PrometheusHook) AfterExecute(task zsched.AnyTask, s *zsched.State) error {
	h.taskCounter.WithLabelValues(task.Name(), string(s.Status)).Inc()
	h.durationHistogram.WithLabelValues(task.Name(), string(s.Status)).Observe(time.Since(s.StartedAt).Seconds())
	return nil
}
