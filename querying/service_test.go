package querying

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testSvc, _ = NewService()
	query1     = &mockQuery{
		res: 1,
	}
	query2 = &mockQuery{
		res: 2,
	}
)

func TestService_DoesUpdateAll(t *testing.T) {
	gauge1 := &mockGauge{}
	gauge2 := &mockGauge{}

	testSvc.Register(query1, gauge1)
	testSvc.Register(query2, gauge2)

	testSvc.UpdateAll()
	if gauge1.i != 1 || gauge2.i != 2 {
		t.Fatal("Gauge was not updated")
	}
}

func TestService_DoesUpdateUponHandler(t *testing.T) {
	gauge1 := &mockGauge{}
	testSvc.Register(query1, gauge1)

	ts := httptest.NewServer(testSvc.Handler(http.NotFoundHandler()))
	defer ts.Close()

	http.Get(ts.URL)

	if gauge1.i != 1 {
		t.Fatal("Gauge was not updated")
	}
}

func TestService_QueryError(t *testing.T) {
	gauge1 := &mockGauge{}
	query3 := &mockQuery{
		err: errors.New("error"),
	}
	testSvc.Register(query3, gauge1)

	err := testSvc.UpdateAll()
	if err == nil {
		t.Fatal("Was expecting error")
	}
	if gauge1.i == 1 {
		t.Fatal("Expecting gauge not to be updated")
	}
}

type mockGauge struct {
	i int
}

func (m *mockGauge) Update(i int) {
	m.i = i
}

type mockQuery struct {
	res int
	err error
}

func (q *mockQuery) Query() (int, error) {
	return q.res, q.err
}
