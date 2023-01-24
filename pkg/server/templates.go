package server

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/Bengissimo/Citybike/pkg/citybike"
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

func applyTemplate(templateBody string, w http.ResponseWriter, data any) error {
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
			array := []int{}
			for i := value - 5; i < value+5; i++ {
				if i >= 0 {
					array = append(array, i)
				}
			}
			return array
		},
	}

	t, err := template.New("template.html").Funcs(funcs).Parse(templateBody)
	if err != nil {
		return err
	}

	err = t.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}
