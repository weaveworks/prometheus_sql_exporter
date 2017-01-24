package config

import (
	"github.com/weaveworks/prometheus_sql_exporter/db"
	"github.com/weaveworks/prometheus_sql_exporter/monitoring"
	"github.com/weaveworks/prometheus_sql_exporter/querying"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Represent the prometheus metric types and the queries to be performed
type Config struct {
	Gauges []Gauge
}
type Gauge struct {
	monitoring.GaugeConfig `yaml:",inline"`
	Queries                []Query
}

type Query struct {
	Name  string
	Query string
}

type Configuration interface {
	RegisterGauges(repo db.Repository, svc querying.Service) error
}

func NewConfiguration(configPath string) (Configuration, error) {
	var cfg Config
	queryBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return &cfg, err
	}
	err = yaml.Unmarshal(queryBytes, &cfg)
	if err != nil {
		return &cfg, err
	}
	return &cfg, err
}

func (c *Config) RegisterGauges(repo db.Repository, svc querying.Service) error {
	for _, pg := range c.Gauges {
		gauge, err := monitoring.NewGauge(pg.GaugeConfig)
		if err != nil {
			return err
		}

		for _, q := range pg.Queries {
			query, err := db.NewIntQuery(db.QueryConfig{
				Repository: repo,
				Query:      q.Query,
			})
			if err != nil {
				return err
			}
			ng, err := monitoring.NewNamedGauge(monitoring.NamedGaugeConfig{
				Gauge: gauge,
				Name:  q.Name,
			})
			if err != nil {
				return err
			}
			svc.Register(query, ng)
		}
	}
	return nil
}
