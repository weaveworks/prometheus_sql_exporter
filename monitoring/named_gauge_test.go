package monitoring

import "testing"

func TestNamedGauge_Update(t *testing.T) {
	gauge := &mockProseGauge{}
	g, _ := NewNamedGauge(NamedGaugeConfig{
		Gauge: gauge,
		Name: "test",
	})
	g.Update(1)
	if gauge.i != 1 || gauge.name != "test" {
		t.Fatal("Gauge did not update")
	}
}

type mockProseGauge struct {
	name string
	i int
}

func (g *mockProseGauge) UpdateInt(name string, i int) {
	g.name = name
	g.i = i
}

