package querying

import (
	"github.com/weaveworks/prometheus_sql_exporter/db"
	"github.com/weaveworks/prometheus_sql_exporter/monitoring"
	"net/http"
)

// Service - Registers queries against gauges
// Once registered a HTTP handler middleware is exposed to update all query/gauge tuples.
type Service interface {
	Register(q db.IntQuery, g monitoring.NamedGauge)
	UpdateAll() error
	Handler(h http.Handler) http.Handler
}

// NewService - Create a new query service
func NewService() (Service, error) {
	return &svc{
		registered: make(map[db.IntQuery]monitoring.NamedGauge),
	}, nil
}

type svc struct {
	registered map[db.IntQuery]monitoring.NamedGauge
}

func (s *svc) UpdateAll() error {
	for q, g := range s.registered {
		count, err := q.Query()
		if err != nil {
			return err
		}
		g.Update(count)
	}
	return nil
}

func (s *svc) Register(q db.IntQuery, g monitoring.NamedGauge) {
	s.registered[q] = g
}

func (s *svc) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := s.UpdateAll()
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		h.ServeHTTP(w, r)
	})
}
