package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusMetrics struct {
	Counter      *prometheus.CounterVec
	GaugeCounter *prometheus.GaugeVec
}

func NewPrometheusMetrics() (*PrometheusMetrics, error) {
	metrics := &PrometheusMetrics{
		Counter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "resume_views_total",
			Help: "Total number of resume views",
		}, []string{"resume_id"}),
		GaugeCounter: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "resume_views_total_gauge", // I am not sure if I need to use gauge metric because resume views never decrease
			Help: "Total number of resume views gauge",
		}, []string{"resume_id"}),
	}

	err := prometheus.Register(metrics.Counter)
	if err != nil {
		return nil, fmt.Errorf("error registering counter metrics: %v", err)
	}

	err = prometheus.Register(metrics.GaugeCounter)
	if err != nil {
		return nil, fmt.Errorf("error registering gauge metrics: %v", err)
	}

	return metrics, nil
}

func (metrics *PrometheusMetrics) Inc(resumeID string) {
	metrics.Counter.WithLabelValues(resumeID).Inc()
	metrics.GaugeCounter.WithLabelValues(resumeID).Inc()
}

func (metrics *PrometheusMetrics) Dec(resumeID string) {
	metrics.GaugeCounter.WithLabelValues(resumeID).Dec()
}
