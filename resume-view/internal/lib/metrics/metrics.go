package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusMetrics struct {
	TotalViewCounter *prometheus.CounterVec
}

func NewPrometheusMetrics() (*PrometheusMetrics, error) {
	metrics := &PrometheusMetrics{
		TotalViewCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "resume_views_total",
			Help: "Total number of resume views",
		}, []string{"resume_id"}),
	}

	err := prometheus.Register(metrics.TotalViewCounter)
	if err != nil {
		return nil, fmt.Errorf("error registering counter metrics: %v", err)
	}

	return metrics, nil
}

func (metrics *PrometheusMetrics) Inc(resumeID string) {
	metrics.TotalViewCounter.WithLabelValues(resumeID).Inc()
}
