// Package citybike creates sqlite database
// and includes database related methods and functions
package citybike

import (
	"database/sql"

	csvtag "github.com/artonge/go-csv-tag/v2"
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

// Journey struct with csvtags
type Journey struct {
	ID           int
	Departure    string `csv:"Departure"`
	Return       string `csv:"Return"`
	DepartureID  int    `csv:"Departure station id"`
	Departure_FI string
	Departure_SE string
	ReturnID     int `csv:"Return station id"`
	Return_FI    string
	Return_SE    string
	Distance     float64 `csv:"Covered distance (m)"`
	Duration     float64 `csv:"Duration (sec.)"`
}

// Station struct with csvtags
type Station struct {
	ID         int    `csv:"ID"`
	FI_Name    string `csv:"Nimi"`
	SE_Name    string `csv:"Namn"`
	FI_Address string `csv:"Osoite"`
	SE_Address string `csv:"Adress"`
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

func (citybike *Citybike) Close() {
	citybike.db.Close()
}

// exec executes the create table query given as string argument
func (citybike *Citybike) exec(createTable string) error {
	_, err := citybike.db.Exec(createTable)
	return err
}

// Download downloads the csv data using csvtag package
func (citybike *Citybike) downloadJourney() ([]Journey, error) {
	files := []string{
		"/Users/bengisu/Downloads/2021-05.csv",
		//"/Users/bengisu/Downloads/2021-06.csv",
		//"/Users/bengisu/Downloads/2021-07.csv",
	}
	journeyTab := []Journey{} // the slice to put the content into
	for _, filename := range files {
		err := csvtag.LoadFromPath(filename, &journeyTab, csvtag.CsvOptions{})
		if err != nil {
			return nil, err
		}
	}
	return journeyTab, nil
}

func (citybike *Citybike) loadJourneyData() error {
	if err := citybike.exec(createJourneyTable); err != nil {
		return err
	}

	journeytab, err := citybike.downloadJourney()
	if err != nil {
		return err
	}

	stmt, err := citybike.db.Prepare(inserttoJourneyTable)
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
	return nil
}

func (citybike *Citybike) downloadStation() ([]Station, error) {

	file := "/Users/bengisu/Downloads/Helsingin_ja_Espoon_kaupunkipyöräasemat_avoin.csv"
	stationTab := []Station{}

	err := csvtag.LoadFromPath(file, &stationTab, csvtag.CsvOptions{})
	if err != nil {
		return nil, err
	}

	return stationTab, nil
}

func (citybike *Citybike) loadStationData() error {
	if err := citybike.exec(createStationTable); err != nil {
		return err
	}

	stationtab, err := citybike.downloadStation()
	if err != nil {
		return err
	}

	stmt, err := citybike.db.Prepare(inserttoStationTable)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, station := range stationtab {
		_, err := stmt.Exec(station.ID, station.FI_Name, station.SE_Name, station.FI_Address, station.SE_Address)
		if err != nil {
			return err
		}
	}

	return nil
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
