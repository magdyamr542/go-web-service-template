package metrics

import (
	"github.com/magdyamr542/go-web-service-template/pkg/version"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	Registry prometheus.Registerer
	// Tracks the count of tags being requested by the GET /resources?tags=[...] endpoint.
	TagsMetric    *prometheus.CounterVec
	VersionMetric *prometheus.GaugeVec
}

func New() Metrics {

	registry := prometheus.DefaultRegisterer

	m := Metrics{
		TagsMetric: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:      "tags_get_count",
			Namespace: "app",
			Subsystem: "handler",
			Help:      "Tracks the count of tags being requested by the GET /resources?tags=[...] endpoint.",
		}, []string{"tags"}),
		VersionMetric: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "info",
			Help: "Information about the app",
		}, []string{"version", "sha"}),
		Registry: registry,
	}

	registry.MustRegister(m.TagsMetric, m.VersionMetric)
	m.VersionMetric.With(prometheus.Labels{"version": version.Version, "sha": version.CommitSHA}).Set(1)

	return m
}
