package monitoring

import (
	"fmt"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

// ProseGauge - A domain model of a prometheus gauge
type ProseGauge interface {
	UpdateInt(name string, i int)
}

type gauge struct {
	pg *prometheus.Gauge
	l  string
}

// ProseGaugeConfig - config
type ProseGaugeConfig struct {
	Namespace string
	Subsystem string
	Name      string
	Label     string
}

// NewProseGauge - Create a new gauge
func NewProseGauge(c ProseGaugeConfig) (ProseGauge, error) {
	g := prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Namespace: c.Namespace,
		Subsystem: c.Subsystem,
		Name:      c.Name,
		Help:      fmt.Sprintf("Prose Gauge for %v", c.Name),
	}, []string{c.Label})
	return &gauge{
		pg: g,
		l:  c.Label,
	}, nil
}

func (g *gauge) UpdateInt(name string, i int) {
	g.pg.With(g.l, name).Set(float64(i))
}
