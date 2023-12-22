package metrics

import (
	"github.com/magdyamr542/go-web-service-template/pkg/version"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	Registry      prometheus.Registerer
	VersionMetric *prometheus.GaugeVec
}

func New() Metrics {

	registry := prometheus.DefaultRegisterer

	m := Metrics{
		VersionMetric: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "info",
			Help: "Information about the app",
		}, []string{"version", "sha"}),
		Registry: registry,
	}

	registry.MustRegister(m.VersionMetric)
	m.VersionMetric.With(prometheus.Labels{"version": version.Version, "sha": version.CommitSHA}).Set(1)

	return m
}
