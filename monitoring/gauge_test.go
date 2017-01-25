package monitoring

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"strings"
)

func TestGauge_UpdateInt(t *testing.T) {
	cfg := ProseGaugeConfig{
		Name: "name",
		Label: "label",
		Namespace: "namespace",
		Subsystem: "subsystem",
	}
	g, _ := NewProseGauge(cfg)
	g.UpdateInt("state", 1)

	ts := httptest.NewServer(promhttp.Handler())
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	respB, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(respB), `namespace_subsystem_name{label="state"} 1`) {
		t.Fatalf("Response should have contained metric, but didn't: %v", string(respB))
	}

}
