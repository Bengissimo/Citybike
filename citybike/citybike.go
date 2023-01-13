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
	Duration INTEGER)`
	createStationTable string = `DROP TABLE IF EXISTS StationList;
	CREATE TABLE IF NOT EXISTS StationList (
	id INTEGER,
	FI_Name TEXT,
	SE_Name TEXT,
	FI_Address TEXT,
	SE_Address TEXT)`
)

// Citybike is a representation of sql.DB
type Citybike struct {
	db *sql.DB
}

// New opens a database using sql.Open assigns this database to the returned Citybike struct
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

// exec executes the create table query given as string argument
func (citybike *Citybike) exec(createTable string) error {
	_, err := citybike.db.Exec(createTable)
	return err
}

// CreateTables creates Journey and Station tables in the Citybike database
func (citybike *Citybike) CreateTables() error {
	if err := citybike.exec(createJourneyTable); err != nil {
		return err
	}

	if err := citybike.exec(createStationTable); err != nil {
		return err
	}
	return nil
}
