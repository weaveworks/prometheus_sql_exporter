package db

// IntQuery - Performs a DB query that returns an integer
type IntQuery interface {
	Query() (int, error)
}

type query struct {
	db Repository
	q  string
}

// QueryConfig - config
type QueryConfig struct {
	Repository Repository
	Query      string
}

// NewIntQuery - constructor
func NewIntQuery(c QueryConfig) (IntQuery, error) {
	return &query{
		q:  c.Query,
		db: c.Repository,
	}, nil
}

func (q *query) Query() (int, error) {
	return q.db.QueryInt(q.q)
}
