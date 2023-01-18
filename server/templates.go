package server

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/Bengissimo/Citybike/citybike"
)

type stationTemplate struct {
	Stations []citybike.Station
}

type journeyTemplate struct {
	Journeys []citybike.Journey
}

func applyTemplate(name string, w http.ResponseWriter, data any) error {
	_, subname, _ := strings.Cut(name, "/")
	floatFunc := template.FuncMap{
		"FormatFloat": func(value float64) string {
			return fmt.Sprintf("%.1f", value)
		},
	}

	t, err := template.New(subname).Funcs(floatFunc).ParseFiles(name)
	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(w, subname, data)
	if err != nil {
		return err
	}

	return nil
}
