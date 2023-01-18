package server

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

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

	st := stationTemplate{
		Stations: stations,
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

	jt := journeyTemplate{
		Journeys: journeys,
	}

	err = applyTemplate("server/journeys.html", w, jt)
	if err != nil {
		fmt.Println(err)
	}
}
