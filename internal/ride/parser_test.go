package ride

import (
	"errors"
	"testing"

	"github.com/kanzafirov/fareestimation/internal/ride/model"
)

type ParseRowResult struct {
	row            []string
	expectedID     string
	expectedSignal *model.LocationSignal
	expectedError  error
}

var distanceResults = []ParseRowResult{
	{[]string{"1", "37.964168", "23.726123", "1405595110"}, "1", &model.LocationSignal{Latitude: 37.964168, Longitude: 23.726123, Timestamp: 1405595110}, nil},
	{[]string{"testid", "35.355555", "35.355555", "1588452415"}, "testid", &model.LocationSignal{Latitude: 35.355555, Longitude: 35.355555, Timestamp: 1588452415}, nil},
	{[]string{"1", "37.964168", "23.726123"}, "", nil, errors.New("")},
	{[]string{"1", "37.964168", "23.726123", "abcd"}, "", nil, errors.New("")},
	{nil, "", nil, errors.New("")},
}

func TestParseRow(t *testing.T) {
	for _, test := range distanceResults {
		id, result, err := parseRow(test.row)

		if id != test.expectedID {
			t.Errorf("[TestParseRow] Failed: expectedId: %v, actualId: %v", test.expectedID, id)
		}

		if result != nil && test.expectedSignal != nil && *result != *test.expectedSignal {
			t.Errorf("[TestParseRow] Failed: expectedSignal: %v, actualSignal: %v", *test.expectedSignal, *result)
		}

		if (err == nil && test.expectedError != nil) || (err != nil && test.expectedError == nil) {
			t.Errorf("[TestParseRow] Failed: expectedError: %v, actualError: %v", test.expectedError, err)
		}
	}
}
