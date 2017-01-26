package config

import (
	"github.com/weaveworks/prometheus_sql_exporter/db"
	"github.com/weaveworks/prometheus_sql_exporter/monitoring"
	"github.com/weaveworks/prometheus_sql_exporter/querying"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Represent the prometheus metric types and the queries to be performed
type proseConfig struct {
	Gauges []Gauge
}

// Gauge - a specific gauge in the configuration
type Gauge struct {
	monitoring.ProseGaugeConfig `yaml:",inline"`
	Queries                     []Query
}

// Query - a specific query in the configuration
type Query struct {
	Name  string
	Query string
}

// ProseConfiguration - interface to main configuration
type ProseConfiguration interface {
	RegisterGauges(repo db.Repository, svc querying.Service) error
}

// NewProseConfiguration - new query/gauge specification from a file
func NewProseConfiguration(configPath string) (ProseConfiguration, error) {
	var cfg proseConfig
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

// RegisterGauges - using the configuration, create new gauges and register them in the service
func (c *proseConfig) RegisterGauges(repo db.Repository, svc querying.Service) error {
	for _, pg := range c.Gauges {
		gauge, err := monitoring.NewProseGauge(pg.ProseGaugeConfig)
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
