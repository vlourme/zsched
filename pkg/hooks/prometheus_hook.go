package hooks

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vlourme/zsched/pkg/state"
	"github.com/vlourme/zsched/pkg/task"
)

type PrometheusHook struct {
	taskCounter       *prometheus.CounterVec
	durationHistogram *prometheus.HistogramVec
}

func NewPrometheusHook() *PrometheusHook {
	taskCounter := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "scheduler_tasks_total",
		Help: "Total number of tasks executed",
	}, []string{"task_name", "status"})

	durationHistogram := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "scheduler_task_duration_seconds",
		Help: "Duration of tasks in seconds",
	}, []string{"task_name", "status"})

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	return &PrometheusHook{
		taskCounter:       taskCounter,
		durationHistogram: durationHistogram,
	}
}

func (h *PrometheusHook) BeforeExecute(task *task.Task, s *state.State) error {
	return nil
}

func (h *PrometheusHook) AfterExecute(task *task.Task, s *state.State) error {
	h.taskCounter.WithLabelValues(task.Name(), string(s.Status)).Inc()
	h.durationHistogram.WithLabelValues(task.Name(), string(s.Status)).Observe(time.Since(s.StartedAt).Seconds())
	return nil
}
