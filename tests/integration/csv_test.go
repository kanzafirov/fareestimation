package ride

import (
	"strings"
	"testing"

	"github.com/kanzafirov/fareestimation/internal/fs/csv"
)

var testData = []string{
	"1,37.966660,23.728308,1405594957",
	"1,37.966627,23.728263,1405594966",
	"1,37.966625,23.728263,1405594974",
}

func TestNewReadChannel(t *testing.T) {
	readChannel, _, err := csv.NewReadChannel("../testdata/minorset.csv")
	if err != nil {
		t.Errorf("[TestNewReadChannel] Failed to get new read channel; err: %v", err)
	}

	var rows []string

	for row := range readChannel {
		rows = append(rows, strings.Join(row, ","))
	}

	if !compare(rows, testData) {
		t.Errorf("[TestNewReadChannel] String arrays not matching; expected: %v; actual: %v", testData, rows)
	}
}

func compare(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
