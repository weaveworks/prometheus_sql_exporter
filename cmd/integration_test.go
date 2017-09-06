// +build integration

package cmd

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// This test expects the postgres db, built from /mocks/Dockerfile.integration-db, to be running on localhost:15432
// We use port 15432 because apparently circle already has something running on 5432
func TestIntegration_SimpleQuery(t *testing.T) {
	file := createTestQueryFile("ns")
	defer os.Remove(file.Name())
	viper.Set(databaseSourceParam, "postgres://postgres@localhost:15432/test?sslmode=disable")
	viper.Set(queriesParam, pathTo(file))
	qSvc := wireUpDomain(log.NewLogfmtLogger(os.Stderr))

	var httpMiddleware http.Handler
	{
		httpMiddleware = promhttp.Handler()
	}

	ts := httptest.NewServer(httpMiddleware)
	defer ts.Close()

	if err := qSvc.UpdateAll(); err != nil {
		t.Error(err)
	}

	resp, err := http.Get(ts.URL + "/metrics")
	if err != nil {
		t.Fatal(err)
	}
	res, _ := ioutil.ReadAll(resp.Body)
	t.Log(fmt.Sprintf("%v: %s", resp.StatusCode, res))
	if !strings.Contains(string(res), `ns_sub_name{test="scheduled"} 1`) {
		t.Fatal("/metrics did not contain expected query")
	}
}

func createTestQueryFile(namespace string) *os.File {
	file, _ := ioutil.TempFile(os.TempDir(), "temp")
	file.Write([]byte(`gauges:
- gauge:
  namespace: "` + namespace + `"
  subsystem: "sub"
  name: "name"
  label: "test"
  queries:
  - name: "scheduled"
    query: "SELECT count(1) FROM data"`))
	return file
}

func pathTo(file *os.File) string {
	thepath, _ := filepath.Abs(file.Name())
	return thepath
}
