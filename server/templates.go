package server

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/Bengissimo/Citybike/citybike"
)

type stationTemplate struct {
	Stations    []citybike.Station
	CurrentPage int
	TotalPages  int
}

type journeyTemplate struct {
	Journeys    []citybike.Journey
	CurrentPage int
	TotalPages  int
}

type singleStationTemplate struct {
	Station citybike.Station
	StartFrom int
	EndingAt int
}

func applyTemplate(name string, w http.ResponseWriter, data any) error {
	_, subname, _ := strings.Cut(name, "/")
	funcs := template.FuncMap{
		"FormatFloat": func(value float64) string {
			return fmt.Sprintf("%.2f", value)
		},
		"MinusOne": func(value int) int {
			return value - 1
		},
		"PlusOne": func(value int) int {
			return value + 1
		},
		"PageRange": func(value int) []int {
			a := []int{}
			for i := value - 5; i < value+5; i++ {
				if i >= 0 {
					a = append(a, i)
				}
			}
			return a
		},
	}

	t, err := template.New(subname).Funcs(funcs).ParseFiles(name)
	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(w, subname, data)
	if err != nil {
		return err
	}

	return nil
}
