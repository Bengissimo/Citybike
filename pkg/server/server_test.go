package server

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/Bengissimo/Citybike/pkg/citybike"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	cb := citybike.NewWithDb(db)

	s := New(cb)

	handler := http.HandlerFunc(s.indexHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := templateIndex
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestStationHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/stations", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id",
		"FI_Name", "SE_Name",
		"FI_Address", "SE_Address",
	}).AddRow(1, "finame", "sename", "fiaddress", "seaddress")
	mock.ExpectPrepare(regexp.QuoteMeta(citybike.SelectStations)).
		ExpectQuery().WithArgs(citybike.PerPage, 0).WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"count(*)"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta(citybike.Stations)).WillReturnRows(rows)

	cb := citybike.NewWithDb(db)

	s := New(cb)
	handler := http.HandlerFunc(s.stationHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Osoite"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestJourneyHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/journeys", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"Departure_ID", "Departure.FI_Name", "Departure.SE_Name",
		"Return_ID", "Return.FI_Name", "Return.SE_Name",
		"Distance", "Duration"}).
		AddRow(1, "1_FIname", "1_SEname", 2, "2_FIname", "2_SEname", 2.00, 0.50)
	mock.ExpectPrepare(regexp.QuoteMeta(citybike.SelectJourneys)).
		ExpectQuery().WithArgs(citybike.PerPage, 0).WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"count(*)"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta(citybike.Journeys)).WillReturnRows(rows)

	cb := citybike.NewWithDb(db)

	s := New(cb)
	handler := http.HandlerFunc(s.journeyHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Departure ID"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSingleViewHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/station/123", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"FI_Name", "SE_Name",
		"FI_Address", "SE_Address",
	}).AddRow("finame", "sename", "fiaddress", "seaddress")
	mock.ExpectPrepare(regexp.QuoteMeta(citybike.SelectSingle)).
		ExpectQuery().WithArgs(123).WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"count(*)"}).AddRow(1)
	mock.ExpectPrepare(regexp.QuoteMeta(citybike.StartingFrom)).
		ExpectQuery().WithArgs(123).WillReturnRows(rows)

	rows = sqlmock.NewRows([]string{"count(*)"}).AddRow(1)
	mock.ExpectPrepare(regexp.QuoteMeta(citybike.EndingAt)).
		ExpectQuery().WithArgs(123).WillReturnRows(rows)

	cb := citybike.NewWithDb(db)

	s := New(cb)
	handler := http.HandlerFunc(s.singleViewHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Name"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
