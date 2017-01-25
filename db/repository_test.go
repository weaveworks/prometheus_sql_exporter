package db

import (
	"testing"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestRepository_CreatePostgres(t *testing.T) {
	db, err := NewDatabase("postgres://localhost")
	if err != nil {
		t.Fatal(err)
	}
	if db == nil {
		t.Fatal("Db is nil")
	}
	// Don't expect to be able to ping it. It doesn't exist.
}

func TestRepository_InvalidURL(t *testing.T) {
	_, err := NewDatabase("bo&gus:\\lala*&^%$/")
	if err == nil {
		t.Fatal("Was expecting error")
	}
}

func TestRepository_ErrorOnUnsupportedDb(t *testing.T) {
	_, err := NewDatabase("notsupported://localhost")
	if err == nil {
		t.Fatal("Should have errored")
	}
}

func TestRepository_QueryInt(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(1)
	mock.ExpectQuery("^SELECT (.+)").WillReturnRows(rows)

	repo := NewRepository(db)

	res, err := repo.QueryInt("SELECT count(1) FROM jobs")
	if err != nil {
		t.Fatal(err)
	}
	if res != 1 {
		t.Fatalf("Expected result of 1, but got: %v", res)
	}
}

// The query expects a single result. So it should error if there are several
func TestRepository_ExpectsSingleArgumentsInResult(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "anothercolumn", "metric"}).
		AddRow("dsa", "aaa", 43)
	mock.ExpectQuery("^SELECT (.+)").WillReturnRows(rows)

	repo := NewRepository(db)

	_, err = repo.QueryInt("SELECT * FROM jobs")
	if err == nil {
		t.Fatal("Was expecting error")
	}
}