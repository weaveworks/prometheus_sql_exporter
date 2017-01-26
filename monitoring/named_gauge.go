package monitoring

// NamedGauge - Maps and updates a ProseGauge with a name
type NamedGauge interface {
	Update(val int)
}

// NamedGaugeConfig - config
type NamedGaugeConfig struct {
	Gauge ProseGauge
	Name  string
}

// NewNamedGauge - constructor
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
