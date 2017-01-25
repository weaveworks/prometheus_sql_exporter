package db

import (
	"database/sql"
	_ "github.com/lib/pq" // For postgres/AWS RDS support. URLs prefixed with "postgres://"
	"net/url"
)

type Repository interface {
	QueryInt(q string) (int, error)
}

type repository struct {
	db *sql.DB
}

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

func NewDatabase(databaseURL string) (conn *sql.DB, err error) {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return
	}
	driver := u.Scheme
	conn, err = sql.Open(driver, databaseURL)
	return
}
