package monitoring

type NamedGauge interface {
	Update(val int)
}

type NamedGaugeConfig struct {
	Gauge ProseGauge
	Name  string
}

func NewNamedGauge(c NamedGaugeConfig) (NamedGauge, error) {
	return &namedGauge{
		pg: c.Gauge,
		n:  c.Name,
	}, nil
}

type namedGauge struct {
	pg ProseGauge
	n  string
}

func (g *namedGauge) Update(val int) {
	g.pg.UpdateInt(g.n, val)
}
