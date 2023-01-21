package citybike

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	csvtag "github.com/artonge/go-csv-tag/v2"
)

var (
	minJourneyDuration time.Duration = 10 * time.Second
	minDistance        float64       = 10
	timeLayout         string        = "2006-01-02T15:04:05"
	stationURLs                      = []string{
		"https://opendata.arcgis.com/datasets/726277c507ef4914b0aec3cbcfcbfafc_0.csv"}
	journeyURLs = []string{
		"https://dev.hsl.fi/citybikes/od-trips-2021/2021-05.csv",
		"https://dev.hsl.fi/citybikes/od-trips-2021/2021-06.csv",
		"https://dev.hsl.fi/citybikes/od-trips-2021/2021-07.csv",
	}
)

// Journey struct with csvtags
type Journey struct {
	ID            int
	DepartureTime string `csv:"Departure"`
	ReturnTime    string `csv:"Return"`
	DepartureID   int    `csv:"Departure station id"`
	DepartureFI   string
	DepartureSE   string
	ReturnID      int `csv:"Return station id"`
	ReturnFI      string
	ReturnSE      string
	Distance      float64 `csv:"Covered distance (m)"`
	Duration      float64 `csv:"Duration (sec.)"`
}

// Station struct with csvtags
type Station struct {
	ID        int    `csv:"ID"`
	NameFI    string `csv:"Nimi"`
	NameSE    string `csv:"Namn"`
	AddressFI string `csv:"Osoite"`
	AddressSE string `csv:"Adress"`
}

// readJourneyCSV reads csv files from specified URLs and returns a slice of Journey struct.
// To load csv it uses csvtag package
func readJourneyCSV() ([]Journey, error) {
	journeyTab := []Journey{}

	brs, err := getURL(journeyURLs)
	if err != nil {
		return nil, err
	}

	for _, br := range brs {

		if err := bomRemover(br); err != nil {
			return nil, err
		}

		if err = csvtag.LoadFromReader(br, &journeyTab); err != nil {
			return nil, err
		}
	}

	validJourneyTab, err := validateJourney(journeyTab)
	if err != nil {
		return nil, err
	}

	return validJourneyTab, nil
}

// getURL reads the response body and returns a slice of bufio.Reader pointers.
func getURL(urls []string) ([]*bufio.Reader, error) {
	brs := []*bufio.Reader{}

	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		resp.Body.Close()
		br := bufio.NewReader(bytes.NewReader(buf))
		brs = append(brs, br)
	}

	return brs, nil
}

// bomRemover removes any Byte Order Mark (BOM) that might be present at the beginning of the buffer.
func bomRemover(br *bufio.Reader) (error) {
	rune, _, err := br.ReadRune()
	if err != nil {
		return err
	}

	if rune != '\uFEFF' {
		br.UnreadRune()
	}

	return nil
}

// validateJourney checks if the journey duration is valid.
// Appends only valid durations to the new validJourneyTab slice.
// Converts distance to kilometers and duration to seconds.
func validateJourney(journeyTab []Journey) ([]Journey, error) {
	validJourneyTab := []Journey{}

	for _, journey := range journeyTab {
		if journey.ReturnTime == "" || journey.DepartureTime == "" {
			continue
		}
		isValid, err := isDurationValid(journey.DepartureTime, journey.ReturnTime)
		if err != nil {
			return nil, err
		}
		if !isValid || journey.Duration < 10 {
			continue
		}
		if journey.Distance < minDistance {
			continue
		}
		journey.Distance /= 1000 // convert to km
		journey.Duration /= 60   // convert to min
		validJourneyTab = append(validJourneyTab, journey)
	}
	return validJourneyTab, nil
}

// parseTime converts string represantaton of a date to a time object.
// If the given value lacks time data, it adds "T00:00:00".
func parseTime(value string) (time.Time, error) {
	if !strings.ContainsAny(value, "T") {
		value = fmt.Sprintf("%sT00:00:00", value)
	}

	t, err := time.Parse(timeLayout, value)
	if err != nil {
		return t, err
	}

	return t, err
}

func isDurationValid(dep string, ret string) (bool, error) {
	tDep, err := parseTime(dep)
	if err != nil {
		return false, err
	}

	tRet, err := parseTime(ret)
	if err != nil {
		return false, err
	}

	if tRet.Sub(tDep) < minJourneyDuration {
		return false, nil
	}

	return true, nil
}

// readStationCSV reads csv file from specified URLs and returns a slice of Station struct.
// To load csv it uses csvtag package.
func readStationCSV() ([]Station, error) {
	stationTab := []Station{}

	brs, err := getURL(stationURLs)
	if err != nil {
		return nil, err
	}

	for _, br := range brs {

		if err := bomRemover(br); err != nil {
			return nil, err
		}

		if err = csvtag.LoadFromReader(br, &stationTab); err != nil {
			return nil, err
		}
	}

	return stationTab, nil
}
