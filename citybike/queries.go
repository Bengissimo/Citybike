package citybike

import (
	"database/sql"
	"errors"
)

var (
	PerPage        int    = 20
	selectStations string = `SELECT id, FI_Name, SE_Name, FI_Address, SE_Address
							FROM StationList
							ORDER BY id
							LIMIT ? OFFSET ?`
	selectJourneys string = `SELECT 
							Departure_ID,
							Departure.FI_Name,
							Departure.SE_Name,
							Return_ID,
							Return.FI_Name,
							Return.SE_Name,
							Distance,
							Duration
							FROM JourneyList
							INNER JOIN StationList AS Departure ON JourneyList.Departure_ID = Departure.id
							INNER JOIN StationList AS Return ON JourneyList.Return_ID = Return.id
							LIMIT ? OFFSET ?`
	selectSingle string = "SELECT FI_Name, SE_Name, FI_Address, SE_Address FROM StationList WHERE id = ?"
)

func (citybike *Citybike) makeQuery(page int, perpage int, query string) (*sql.Rows, error) {
	stmt, err := citybike.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(perpage, page*perpage)
	if err != nil {
		return nil, err
	}
	return rows, err
}

func (citybike *Citybike) GetJourneyRows(pagenum int) ([]Journey, error) {
	rows, err := citybike.makeQuery(pagenum, PerPage, selectJourneys)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	journeys := []Journey{}
	for rows.Next() {
		theJourney := Journey{}
		rows.Scan(&theJourney.DepartureID,
			&theJourney.DepartureFI,
			&theJourney.DepartureSE,
			&theJourney.ReturnID,
			&theJourney.ReturnFI,
			&theJourney.ReturnSE,
			&theJourney.Distance,
			&theJourney.Duration)
		journeys = append(journeys, theJourney)
	}

	return journeys, nil
}

func (citybike *Citybike) GetStationRows(pagenum int) ([]Station, error) {
	rows, err := citybike.makeQuery(pagenum, PerPage, selectStations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stations := []Station{}
	for rows.Next() {
		theStation := Station{}
		rows.Scan(&theStation.ID,
			&theStation.NameFI,
			&theStation.NameSE,
			&theStation.AddressFI,
			&theStation.AddressSE)
		stations = append(stations, theStation)
	}

	return stations, nil
}

func (citybike *Citybike) GetSingleStation(id int) (*Station, error) {
	stmt, err := citybike.db.Prepare(selectSingle)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("Station not found")
	}

	theStation := &Station{
		ID: id,
	}
	rows.Scan(
		&theStation.NameFI,
		&theStation.NameSE,
		&theStation.AddressFI,
		&theStation.AddressSE)

	return theStation, nil
}
