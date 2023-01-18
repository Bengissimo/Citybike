// Package citybike creates sqlite database
// and includes database related methods and functions
package citybike

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var (
	createJourneyTable string = `DROP TABLE IF EXISTS JourneyList;
	CREATE TABLE IF NOT EXISTS JourneyList (
	id INTEGER PRIMARY KEY,
	Departure_ID INTEGER,
	Return_ID INTEGER,
	Distance FLOAT,
	Duration FLOAT)`
	createStationTable string = `DROP TABLE IF EXISTS StationList;
	CREATE TABLE IF NOT EXISTS StationList (
	id INTEGER,
	FI_Name TEXT,
	SE_Name TEXT,
	FI_Address TEXT,
	SE_Address TEXT)`
	inserttoJourneyTable string = `INSERT INTO JourneyList (
	Departure_ID,
	Return_ID,
	Distance,
	Duration)
	VALUES (?, ?, ?, ?)`
	inserttoStationTable string = `INSERT INTO StationList (
	id,
	FI_Name,
	SE_Name,
	FI_Address,
	SE_Address)
	VALUES (?, ?, ?, ?, ?)`
)

// Citybike is a representation of sql.DB
type Citybike struct {
	db *sql.DB
}

// New returns a Citybike struct containing a sqlite database
func New(path string) (*Citybike, error) {
	sqldb, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	db := &Citybike{
		db: sqldb,
	}
	return db, nil
}

func (citybike *Citybike) Close() {
	citybike.db.Close()
}

func (citybike *Citybike) exec(createTable string) error {
	_, err := citybike.db.Exec(createTable)
	return err
}

// loadJourneyData creates Journey table in the database and inserts data
func (citybike *Citybike) loadJourneyData() error {
	tx, err := citybike.db.Begin()
	if err != nil {
		return err
	}

	if err := citybike.exec(createJourneyTable); err != nil {
		return err
	}

	journeytab, err := readJourneyCSV()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(inserttoJourneyTable)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, journey := range journeytab {
		_, err := stmt.Exec(journey.DepartureID, journey.ReturnID, journey.Distance, journey.Duration)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

// loadStationData creates Station table in the database and inserts data
func (citybike *Citybike) loadStationData() error {
	tx, err := citybike.db.Begin()
	if err != nil {
		return err
	}

	if err := citybike.exec(createStationTable); err != nil {
		return err
	}

	stationtab, err := readStationCSV()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(inserttoStationTable)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, station := range stationtab {
		_, err := stmt.Exec(station.ID, station.NameFI, station.NameSE, station.AddressFI, station.AddressSE)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (citybike *Citybike) LoadData() error {
	if err := citybike.loadJourneyData(); err != nil {
		return err
	}

	if err := citybike.loadStationData(); err != nil {
		return err
	}

	return nil
}
