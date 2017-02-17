package db

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq" // For postgres/AWS RDS support. URLs prefixed with "postgres://"
	_ "github.com/go-sql-driver/mysql" // For MySQL support. URLs prefixed with "mysql://"
	"net/url"
	"strings"
)

// Repository - Perform queries on a db and return a metric
type Repository interface {
	QueryInt(q string) (int, error)
}

type repository struct {
	db *sql.DB
}

// NewRepository - constructor
func NewRepository(DB *sql.DB) Repository {
	return &repository{
		db: DB,
	}
}

func (r *repository) QueryInt(q string) (count int, err error) {
	err = r.db.QueryRow(q).Scan(&count)
	if err != nil {
		return
	}
	return
}

// formatDatabaseDSN - generate database DSN in format expected by the driver from the string passed by the user
// postgres expects postgres:// scheme to be set, while mysql would interpret mysql:// prefix as user definition
func formatDatabaseDSN(driver, databaseURL string) string {
	if driver == "mysql" {
		return fmt.Sprintf(strings.TrimPrefix(databaseURL, fmt.Sprintf("%s://", driver)))
	}
	return databaseURL
}

// NewDatabase - instantiate a DB connection from a string url
func NewDatabase(databaseURL string) (conn *sql.DB, err error) {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return
	}
	driver := u.Scheme
	databaseDSN := formatDatabaseDSN(driver, databaseURL)
	conn, err = sql.Open(driver, databaseDSN)
	return
}
