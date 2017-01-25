package db

import "testing"

func TestIntQuery_Query(t *testing.T) {
	repo := &mockRepository{
		res: 1,
	}
	cfg := QueryConfig{
		Query:      "test",
		Repository: repo,
	}
	q, _ := NewIntQuery(cfg)
	if res, _ := q.Query(); res != 1 {
		t.Fatalf("Expecting mock db to return 1, received %q", res)
	}
}

type mockRepository struct {
	query string
	res   int
	err   error
}

func (r *mockRepository) QueryInt(q string) (int, error) {
	r.query = q
	return r.res, r.err
}
