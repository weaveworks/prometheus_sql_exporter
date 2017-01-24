package monitoring

import (
	"fmt"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

type Gauge interface {
	UpdateInt(name string, i int)
}

type gauge struct {
	pg *prometheus.Gauge
	l  string
}

type GaugeConfig struct {
	Namespace string
	Subsystem string
	Name      string
	Label     string
}

func NewGauge(c GaugeConfig) (Gauge, error) {
	g := prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Namespace: c.Namespace,
		Subsystem: c.Subsystem,
		Name:      c.Name,
		Help:      fmt.Sprintf("Gauge for %v", c.Name),
	}, []string{c.Label})
	return &gauge{
		pg: g,
		l:  c.Label,
	}, nil
}

func (g *gauge) UpdateInt(name string, i int) {
	g.pg.With(g.l, name).Set(float64(i))
}
