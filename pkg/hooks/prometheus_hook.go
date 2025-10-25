package hooks

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vlourme/scheduler/pkg/state"
	"github.com/vlourme/scheduler/pkg/task"
)

type PrometheusHook struct {
	taskCounter *prometheus.CounterVec
}

func NewPrometheusHook() *PrometheusHook {
	taskCounter := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "scheduler_tasks_total",
		Help: "Total number of tasks executed",
	}, []string{"task_name", "status"})

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	return &PrometheusHook{
		taskCounter: taskCounter,
	}
}

func (h *PrometheusHook) BeforeExecute(task *task.Task, s *state.State) error {
	return nil
}

func (h *PrometheusHook) AfterExecute(task *task.Task, s *state.State) error {
	h.taskCounter.WithLabelValues(task.Name(), string(s.Status)).Inc()
	return nil
}
