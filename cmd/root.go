package cmd

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/weaveworks/prometheus_sql_exporter/config"
	"github.com/weaveworks/prometheus_sql_exporter/db"
	"github.com/weaveworks/prometheus_sql_exporter/querying"
	"net/http"
	"os"
)

const (
	// Parameters
	databaseSourceParam = "dbsource"
	listenParam         = "listen"
	queriesParam        = "queries"

	// Defaults
	defaultServerPort     string = ":80"
	defaultDatabaseSource string = ""
	defaultQueries        string = "queries.yaml"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	bindLocalFlag(rootCmd, databaseSourceParam, defaultDatabaseSource, `Database source name; includes the DB driver as the scheme. The default is a temporary, file-based DB`)
	bindLocalFlag(rootCmd, listenParam, defaultServerPort, `Listen address for API clients`)
	bindLocalFlag(rootCmd, queriesParam, defaultQueries, `Path to yaml file which describes metrics and queries`)
	viper.AutomaticEnv() // read in environment variables that match
}

func bindLocalFlag(c *cobra.Command, name string, value string, help string) {
	c.Flags().String(name, value, help)
	viper.BindPFlag(name, c.Flags().Lookup(name))
}

var rootCmd = &cobra.Command{
	Use:   "prose",
	Short: "Monitor a database and expose metrics for prometheus",
	Long:  `This service will monitor a database for specified queries and expose them to prometheus`,
	Run: func(cmd *cobra.Command, args []string) {
		// Logger component.
		var logger log.Logger
		{
			logger = log.NewLogfmtLogger(os.Stderr)
			logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
			logger = log.NewContext(logger).With("caller", log.DefaultCaller)
		}

		qSvc := wireUpDomain(logger)

		// Error channel
		errc := make(chan error)

		var httpMiddleware http.Handler
		{
			httpMiddleware = promhttp.Handler()
			httpMiddleware = qSvc.Handler(httpMiddleware)
		}

		// Start prometheus metrics endpoint
		go func() {
			logger.Log("stage", "httpserver", "addr", viper.GetString(listenParam))
			mux := http.NewServeMux()
			mux.Handle("/metrics", httpMiddleware)
			errc <- http.ListenAndServe(viper.GetString(listenParam), mux)
		}()

		logger.Log("exiting", <-errc)
	},
}

// Execute - run the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func wireUpDomain(logger log.Logger) querying.Service {
	// Create database connection and repository
	database, err := db.NewDatabase(viper.GetString(databaseSourceParam))
	if err != nil {
		logger.Log("stage", "db init", "err", err)
		os.Exit(1)
	}
	repository := db.NewRepository(database)

	// Create querying service
	qSvc, err := querying.NewService()
	if err != nil {
		logger.Log("stage", "query svc init", "err", err)
		os.Exit(1)
	}

	// Register queries and gauges
	cfg, err := config.NewProseConfiguration(viper.GetString(queriesParam))
	if err != nil {
		logger.Log("stage", "configuration", "err", err)
		os.Exit(1)
	}
	cfg.RegisterGauges(repository, qSvc)
	if err != nil {
		logger.Log("stage", "register gauges", "err", err)
		os.Exit(1)
	}

	return qSvc
}
