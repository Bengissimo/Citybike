package citybike

var (
	countStations string = "SELECT count(*) FROM StationList"
	countJourneys string = `SELECT count(*)
							FROM JourneyList
							INNER JOIN StationList AS Departure ON JourneyList.Departure_ID = Departure.id
							INNER JOIN StationList AS Return ON JourneyList.Return_ID = Return.id`
	JourneysStartFrom string = `SELECT count(*) 
								FROM JourneyList
								WHERE Departure_ID = ?`
	JourneysEndingAt string = `SELECT count(*) 
								FROM JourneyList
								WHERE Return_ID = ?`
)

func (citybike *Citybike) CountStations() (int, error) {
	rows, err := citybike.db.Query(countStations)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	rows.Next()
	rows.Scan(&count)
	return count, nil
}

func (citybike *Citybike) CountJourneys() (int, error) {
	rows, err := citybike.db.Query(countJourneys)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	rows.Next()
	rows.Scan(&count)
	return count, nil
}

func (citybike *Citybike) CountStationJourneys(id int, prep string) (int, error) {
	count := 0

	stmt, err := citybike.db.Prepare(prep)
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
