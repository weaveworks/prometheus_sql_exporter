package db

import (
	"database/sql"
	"net/url"
	"github.com/go-kit/kit/log"
)

type Repository interface {
	QueryInt(q string) (int, error)
}

type repository struct {
	db *sql.DB
	Logger log.Logger
}

type RepositoryConfig struct {
	DatabaseUrl string
	Logger log.Logger
}

func NewRepository(c RepositoryConfig) (r Repository, err error) {
	u, err := url.Parse(c.DatabaseUrl)
	if err != nil {
		c.Logger.Log("err", err)
		return
	}
	c.Logger.Log("url", u, "scheme", u.Scheme)
	driver := u.Scheme

	conn, err := sql.Open(driver, c.DatabaseUrl)
	if err != nil {
		c.Logger.Log("err", err)
		return
	}
	r = &repository{
		db: conn,
	}
	return
}

func (r *repository) QueryInt(q string) (count int, err error) {
	err = r.db.QueryRow(q).Scan(&count)
	if err != nil {
		return
	}
	return
}

