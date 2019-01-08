[![CircleCI](https://circleci.com/gh/weaveworks/prometheus_sql_exporter/tree/master.svg?style=svg&circle-token=584f0d5f600891d52e2b5fa7f20a079afd9b47a2)](https://circleci.com/gh/weaveworks/prometheus_sql_exporter/tree/master) [![Go Report Card](https://goreportcard.com/badge/github.com/weaveworks/prometheus_sql_exporter)](https://goreportcard.com/report/github.com/weaveworks/prometheus_sql_exporter) [![Coverage Status](https://coveralls.io/repos/github/weaveworks/prometheus_sql_exporter/badge.svg?branch=master)](https://coveralls.io/github/weaveworks/prometheus_sql_exporter?branch=master)

# prometheus_sql_exporter

The **Pro**metheus **S**QL **E**xporter (PROSE) is a tool that converts user-specified SQL queries into Prometheus metrics. It has the following features

-   Create any number of Prometheus gauges
-   Attach any number of raw SQL queries to a Prometheus gauge
-   Support PostgreSQL and MySQL

Caveats:

-   Only gauges are supported
-   All SQL queries must return a single integer metric (e.g. `count()`)

## Getting started

It is intended that this tool is used as a Docker container.

To use with Kubernetes, use a manifest that [looks like the provided example](./deploy/k8s/prose.yaml)

### CLI Arguments

```
$ docker run -it quay.io/weaveworks/prometheus_sql_exporter --help
This service will monitor a database for specified queries and expose them to prometheus

Usage:
prose [flags]
prose [command]

Available Commands:
version     Output the version of prose

Flags:
    --dbsource string   Database source name; includes the DB driver as the scheme. E.g. postgres://user:password@localhost:5432/database?sslmode=disable or mysql://user:password@tcp(localhost:3306)/database?tls=skip-verify
    --listen string     Listen address for API clients (default ":80")
    --queries string    Path to yaml file which describes metrics and queries (default "queries.yaml")

Use "prose [command] --help" for more information about a command.

```

### Queries configuration

Below is an example queries yaml file:

```yaml
gauges:
- gauge:
  namespace: "mynamespace"
  subsystem: "mysubsystem"
  name: "some_name"
  label: "state"
  queries:
  - name: "number"
    query: "SELECT count(1) FROM database WHERE id > 10"
  - name: "ok"
    query: "SELECT count(1) FROM database"
...
```

-   `gauges` is a list of Prometheus gauges
-   `gauge` is a single instance
-   `namespace`, `subsystem`, and `name` form the name of the Prometheus object
-   `label` is the Prometheus label corresponding to the queries
-   `queries` is a list of queries to perform
-   `queries.name` is a name applied to the label.
-   `queries.query` is the SQL query to perform on the DB.

This configuration will produce two Prometheus metrics, in the form:

```
# HELP mynamespace_mysubsystem_some_name Prose Guage for some_name
# TYPE mynamespace_mysubsystem_some_name gauge
mynamespace_mysubsystem_some_name{state="number"} 90
mynamespace_mysubsystem_some_name{state="ok"} 100
```

## Development

The build lifecycle is controlled by the Makefile. [See the Makefile for details](./Makefile)

### Release

To perform a release, create a new GH release.

## <a name="help"></a>Getting Help

If you have any questions about, feedback for or problems with `prometheus_sql_exporter`:

- Invite yourself to the <a href="https://slack.weave.works/" target="_blank">Weave Users Slack</a>.
- Ask a question on the [#general](https://weave-community.slack.com/messages/general/) slack channel.
- [File an issue](https://github.com/weaveworks/prometheus_sql_exporter/issues/new).

Your feedback is always welcome!
