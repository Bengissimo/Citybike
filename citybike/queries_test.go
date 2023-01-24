package citybike

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCitybike_GetJourneyRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"Departure_ID", "Departure.FI_Name", "Departure.SE_Name",
		"Return_ID", "Return.FI_Name", "Return.SE_Name",
		"Distance", "Duration"}).
		AddRow(1, "1_FIname", "1_SEname", 2, "2_FIname", "2_SEname", 2.0, 0.5)

	mock.ExpectPrepare(regexp.QuoteMeta(SelectJourneys)).
		ExpectQuery().WithArgs(PerPage, 0).WillReturnRows(rows)

	cb := &Citybike{
		db: db,
	}

	j, err := cb.GetJourneyRows(0)
	if err != nil {
		t.Fatalf("error occured %v", err)
	}
	if len(j) != 1 {
		t.Fatalf("size did not match")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

var expStation = Station{
	ID:        1,
	NameFI:    "FIname",
	NameSE:    "SEname",
	AddressFI: "FIaddress",
	AddressSE: "SEaddress",
}

func TestCitybike_GetStationRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"id",
		"FI_Name", "SE_Name",
		"FI_Address", "SE_Address",
	}).AddRow(
		expStation.ID,
		expStation.NameFI, expStation.NameSE,
		expStation.AddressFI, expStation.AddressSE)

	mock.ExpectPrepare(regexp.QuoteMeta(SelectStations)).
		ExpectQuery().WithArgs(PerPage, 0).WillReturnRows(rows)

	cb := &Citybike{
		db: db,
	}

	s, err := cb.GetStationRows(0)
	if err != nil {
		t.Fatalf("error occured %v", err)
	}
	if len(s) != 1 {
		t.Fatalf("size did not match")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCitybike_GetSingleStation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"FI_Name", "SE_Name",
		"FI_Address", "SE_Address",
	}).AddRow(
		expStation.NameFI, expStation.NameSE,
		expStation.AddressFI, expStation.AddressSE)

	mock.ExpectPrepare(regexp.QuoteMeta(SelectSingle)).
		ExpectQuery().WithArgs(expStation.ID).WillReturnRows(rows)

	falseID := 1000
	mock.ExpectPrepare(regexp.QuoteMeta(SelectSingle)).
		ExpectQuery().WithArgs(falseID).WillReturnRows(sqlmock.NewRows([]string{}))

	cb := &Citybike{
		db: db,
	}

	s, err := cb.GetSingleStation(expStation.ID)
	if err != nil {
		t.Fatalf("error occured %v", err)
	}

	if !reflect.DeepEqual(s, &expStation) {
		t.Errorf("s:%v exp:%v", s, expStation)
		t.Fatalf("The station did not match")
	}

	_, err = cb.GetSingleStation(falseID)
	if err == nil || err.Error() != "Station not found" {
		t.Fatalf("Expected error not occured")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
