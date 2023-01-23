package citybike

import (
	"bufio"
	"bytes"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	testCases := []struct {
		value    string
		expected time.Time
	}{
		{"2022-12-25", time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC)},
		{"2022-12-25T12:00:00", time.Date(2022, 12, 25, 12, 0, 0, 0, time.UTC)},
	}
	for _, tc := range testCases {
		result, err := parseTime(tc.value)
		if err != nil {
			t.Errorf("Error parsing time %s: %v", tc.value, err)
		}
		if !result.Equal(tc.expected) {
			t.Errorf("Error: expected %v but got %v", tc.expected, result)
		}
	}
	// testing invalid input
	_, err := parseTime("2022-25-12")
	if err == nil {
		t.Errorf("Error: expected error but got nil")
	}
}

func TestIsDurationValid(t *testing.T) {
	testCases := []struct {
		dep      string
		ret      string
		expected bool
	}{
		{"2022-12-25T12:00:00", "2022-12-25T12:30:00", true},
		{"2022-12-25T12:00:00", "2022-12-25T11:59:00", false},
		{"2022-12-25T12:00:00", "2022-12-25T12:00:09", false},
	}
	minJourneyDuration = time.Minute
	for _, tc := range testCases {
		result, err := isDurationValid(tc.dep, tc.ret)
		if err != nil {
			t.Errorf("Error parsing time %v", err)
		}
		if result != tc.expected {
			t.Errorf("Error: expected %v but got %v", tc.expected, result)
		}
	}
}

func TestBomRemover(t *testing.T) {
	testCases := []struct {
		body     string
		expected string
	}{
		{"\uFEFFhello world", "hello world"},
		{"hello world", "hello world"},
	}
	for _, tc := range testCases {
		br := bufio.NewReader(strings.NewReader(tc.body))
		err := bomRemover(br)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		buf := &bytes.Buffer{}
		_, err = buf.ReadFrom(br)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if !reflect.DeepEqual(buf.String(), tc.expected) {
			t.Errorf("Error: expected %v but got %v", tc.expected, tc.body)
		}
	}
}

func TestValidateJourney(t *testing.T) {
	testCases := []struct {
		journeyTab []Journey
		expected   []Journey
	}{
		{[]Journey{{DepartureTime: "2022-12-25T12:00:00", ReturnTime: "2022-12-25T12:30:00", Distance: 2000, Duration: 30}},
			[]Journey{{DepartureTime: "2022-12-25T12:00:00", ReturnTime: "2022-12-25T12:30:00", Distance: 2, Duration: 0.5}}},
		{[]Journey{{DepartureTime: "2022-12-25T12:00:00", ReturnTime: "2022-12-25T12:30:00", Distance: 2, Duration: 30}},
			[]Journey{}},
		{[]Journey{{DepartureTime: "2022-12-25T12:00:00", ReturnTime: "2022-12-25T11:59:00", Distance: 2000, Duration: 30}},
			[]Journey{}},
		{[]Journey{{DepartureTime: "2022-12-25T12:00:00", ReturnTime: "2022-12-25T12:00:01", Distance: 2000, Duration: 9}},
			[]Journey{}},
		{[]Journey{{DepartureTime: "", ReturnTime: "2022-12-25T12:00:01", Distance: 2000, Duration: 30}},
			[]Journey{}},
	}
	for _, tc := range testCases {
		result, err := validateJourney(tc.journeyTab)
		if err != nil {
			t.Errorf("Error parsing time %v: %v", tc.journeyTab, err)
		}
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("Error: expected %v but got %v", tc.expected, result)
		}
	}
}
