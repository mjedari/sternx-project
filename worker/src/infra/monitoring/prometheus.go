package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Prometheus struct {
	tasksReceived *prometheus.CounterVec
	tasksDone     *prometheus.CounterVec
}

func NewPrometheus() *Prometheus {
	tasksReceived := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tasks_received",
			Help: "Memory usage of the application.",
		},
		[]string{"worker"},
	)

	tasksDone := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tasks_done",
			Help: "Memory usage of the application.",
		},
		[]string{"worker"},
	)

	prometheus.MustRegister(tasksDone, tasksReceived)

	return &Prometheus{tasksReceived: tasksReceived, tasksDone: tasksDone}
}

func (p Prometheus) ReceivedTasks(queue string) {
	p.tasksReceived.WithLabelValues(queue).Inc()
}

func (p Prometheus) DoneTasks(queue string) {
	p.tasksDone.WithLabelValues(queue).Inc()
}
