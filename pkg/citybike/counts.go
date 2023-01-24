package citybike

var (
	Stations string = "SELECT count(*) FROM StationList"
	Journeys string = `SELECT count(*)
							FROM JourneyList
							INNER JOIN StationList AS Departure ON JourneyList.Departure_ID = Departure.id
							INNER JOIN StationList AS Return ON JourneyList.Return_ID = Return.id`
	StartingFrom string = `SELECT count(*) 
								FROM JourneyList
								WHERE Departure_ID = ?`
	EndingAt string = `SELECT count(*) 
								FROM JourneyList
								WHERE Return_ID = ?`
)

func (citybike *Citybike) count(query string) (int, error) {
	rows, err := citybike.db.Query(query)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	if rows.Next() {
		rows.Scan(&count)
	}

	return count, nil
}

func (citybike *Citybike) CountStations() (int, error) {
	return citybike.count(Stations)
}

func (citybike *Citybike) CountJourneys() (int, error) {
	return citybike.count(Journeys)
}

func (citybike *Citybike) countStationJourneys(id int, query string) (int, error) {
	count := 0

	stmt, err := citybike.db.Prepare(query)
	if err != nil {
		return count, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return count, err
	}
	defer rows.Close()

	if !rows.Next() {
		return count, nil
	}
	rows.Scan(&count)

	return count, nil
}

func (citybike *Citybike) CountStartingFrom(id int) (int, error) {
	return citybike.countStationJourneys(id, StartingFrom)
}

func (citybike *Citybike) CountEndingAt(id int) (int, error) {
	return citybike.countStationJourneys(id, EndingAt)
}
