package citybike

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCitybikeCountStations(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	expected := 123

	rows := sqlmock.NewRows([]string{"count(*)"}).
		AddRow(expected)

	mock.ExpectQuery(regexp.QuoteMeta(Stations)).WillReturnRows(rows)

	cb := &Citybike{
		db: db,
	}

	count, err := cb.CountStations()
	if err != nil {
		t.Fatalf("error occured %v", err)
	}
	if count != expected {
		t.Fatalf("Count did not match")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCitybikeCountJourneys(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	expected := 123

	rows := sqlmock.NewRows([]string{"count(*)"}).
		AddRow(expected)

	mock.ExpectQuery(regexp.QuoteMeta(Journeys)).WillReturnRows(rows)

	cb := &Citybike{
		db: db,
	}

	count, err := cb.CountJourneys()
	if err != nil {
		t.Fatalf("error occured %v", err)
	}
	if count != expected {
		t.Fatalf("Count did not match")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCitybikeCountStartingFrom(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	expected := 123

	rows := sqlmock.NewRows([]string{"count(*)"}).
		AddRow(expected)

	mock.ExpectPrepare(regexp.QuoteMeta(StartingFrom)).
		ExpectQuery().WithArgs(1).WillReturnRows(rows)

	cb := &Citybike{
		db: db,
	}

	count, err := cb.CountStartingFrom(1)
	if err != nil {
		t.Fatalf("error occured %v", err)
	}
	if count != expected {
		t.Fatalf("Count did not match")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCitybikeCountEndingAt(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	expected := 123

	rows := sqlmock.NewRows([]string{"count(*)"}).
		AddRow(expected)

	mock.ExpectPrepare(regexp.QuoteMeta(EndingAt)).
		ExpectQuery().WithArgs(1).WillReturnRows(rows)

	cb := &Citybike{
		db: db,
	}

	count, err := cb.CountEndingAt(1)
	if err != nil {
		t.Fatalf("error occured %v", err)
	}
	if count != expected {
		t.Fatalf("Count did not match")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
