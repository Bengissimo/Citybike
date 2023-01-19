package server

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/Bengissimo/Citybike/citybike"
)

type Server struct {
	db *citybike.Citybike
}

func New(cb *citybike.Citybike) *Server {
	srv := &Server{
		db: cb,
	}

	http.HandleFunc("/", srv.indexHandler)
	http.HandleFunc("/journeys", srv.journeyHandler)
	http.HandleFunc("/stations", srv.stationHandler)
	http.HandleFunc("/station/", srv.singleViewHandler)

	return srv
}

func (server *Server) Run() error {
	if err := http.ListenAndServe(":8000", nil); err != nil {
		return err
	}

	return nil
}

func (server *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("server/index.html")
	t.Execute(w, nil)
}

func (server *Server) stationHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("p")
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		pageNum = 0
	}

	stations, err := server.db.GetStationRows(pageNum)
	if err != nil {
		fmt.Println(err)
	}

	count, err := server.db.CountStations()
	if err != nil {
		fmt.Println(err)
	}
	totalPages := math.Ceil(float64(count) / float64(citybike.PerPage))

	st := stationTemplate{
		Stations:    stations,
		CurrentPage: pageNum,
		TotalPages:  int(totalPages) - 1,
	}

	err = applyTemplate("server/stations.html", w, st)
	if err != nil {
		fmt.Println(err)
	}
}

func (server *Server) journeyHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("p")
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		pageNum = 0
	}

	journeys, err := server.db.GetJourneyRows(pageNum)
	if err != nil {
		fmt.Println(err)
	}

	count, err := server.db.CountJourneys()
	if err != nil {
		fmt.Println(err)
	}
	totalPages := math.Ceil(float64(count) / float64(citybike.PerPage))

	jt := journeyTemplate{
		Journeys:    journeys,
		CurrentPage: pageNum,
		TotalPages:  int(totalPages) - 1,
	}

	err = applyTemplate("server/journeys.html", w, jt)
	if err != nil {
		fmt.Println(err)
	}
}

func (server *Server) singleViewHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.Split(r.URL.Path, "/")[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		id = 0
	}
	theStation, err := server.db.GetSingleStation(id)
	if err != nil {
		fmt.Println(err)
	}

	start, end, err := server.db.CountStationJourney(id)
	if err != nil {
		fmt.Println(err)
	}

	sst := singleStationTemplate{
		Station:   *theStation,
		StartFrom: start,
		EndingAt: end,
	}

	err = applyTemplate("server/singleview.html", w, sst)
	if err != nil {
		fmt.Println(err)
	}
}
