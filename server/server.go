package server

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	_ "embed"

	"github.com/Bengissimo/Citybike/citybike"
)

type Server struct {
	cb *citybike.Citybike
	mux *http.ServeMux
}

var (
	//go:embed index.html
	templateIndex string
	//go:embed journeys.html
	templateJourney string
	//go:embed stations.html
	templateStation string
	//go:embed singleview.html
	templateSingleView string
)

func New(cb *citybike.Citybike) *Server {
	srv := &Server{
		cb: cb,
		mux : http.NewServeMux(),
	}

	srv.mux.HandleFunc("/", srv.indexHandler)
	srv.mux.HandleFunc("/journeys", srv.journeyHandler)
	srv.mux.HandleFunc("/stations", srv.stationHandler)
	srv.mux.HandleFunc("/station/", srv.singleViewHandler)

	return srv
}

func (server *Server) Run() error {
	if err := http.ListenAndServe(":8000", server.mux); err != nil {
		return err
	}

	return nil
}

func (server *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("index.html").Parse(templateIndex)
	if err != nil {
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
		log.Println(err)
	}
}

func (server *Server) stationHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("p")
	pageNum := 0
	var err error
	if page != "" {
		if pageNum, err = strconv.Atoi(page); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}

	stations, err := server.cb.GetStationRows(pageNum)
	if err != nil {
		http.Error(w, "Unable to load station data", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	count, err := server.cb.CountStations()
	if err != nil {
		http.Error(w, "Unable to fetch count from database", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	totalPages := math.Ceil(float64(count) / float64(citybike.PerPage))

	if pageNum >= int(totalPages) {
		http.Redirect(w, r, fmt.Sprintf("/stations?p=%d", int(totalPages)-1), http.StatusTemporaryRedirect)
		return
	}

	st := stationTemplate{
		Stations:    stations,
		CurrentPage: pageNum,
		TotalPages:  int(totalPages) - 1,
	}

	err = applyTemplate(templateStation, w, st)
	if err != nil {
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
		log.Println(err)
	}
}

func (server *Server) journeyHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("p")
	pageNum := 0
	var err error
	if page != "" {
		if pageNum, err = strconv.Atoi(page); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}

	journeys, err := server.cb.GetJourneyRows(pageNum)
	if err != nil {
		http.Error(w, "Unable to load journey data", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	count, err := server.cb.CountJourneys()
	if err != nil {
		http.Error(w, "Unable to fetch count from database", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	totalPages := math.Ceil(float64(count) / float64(citybike.PerPage))

	if pageNum >= int(totalPages) {
		http.Redirect(w, r, fmt.Sprintf("/journeys?p=%d", int(totalPages)-1), http.StatusTemporaryRedirect)
		return
	}

	jt := journeyTemplate{
		Journeys:    journeys,
		CurrentPage: pageNum,
		TotalPages:  int(totalPages) - 1,
	}

	err = applyTemplate(templateJourney, w, jt)
	if err != nil {
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
		log.Println(err)
	}
}

func (server *Server) singleViewHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.Split(r.URL.Path, "/")[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	theStation, err := server.cb.GetSingleStation(id)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	start, err := server.cb.CountStartingFrom(id)
	if err != nil {
		log.Println(err)
	}

	end, err := server.cb.CountEndingAt(id)
	if err != nil {
		log.Println(err)
	}

	sst := singleStationTemplate{
		Station:   *theStation,
		StartFrom: start,
		EndingAt:  end,
	}

	err = applyTemplate(templateSingleView, w, sst)
	if err != nil {
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
		log.Println(err)
	}
}
