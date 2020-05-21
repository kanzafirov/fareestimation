package ride

import (
	"fmt"
	"strconv"

	"github.com/kanzafirov/fareestimation/internal/ride/model"
)

const buffSize = 50

// NewParseChannel gets an input channel as parameter, passing string arrays; for each of them
// a Ride model is generated and returns a new model.Ride channel, where the parsed data is sent
func NewParseChannel(in <-chan []string) (<-chan model.Ride, <-chan error) {
	out := make(chan model.Ride, buffSize)
	outerror := make(chan error)

	go func() {
		var id string
		var locationsInRide []model.LocationSignal

		for row := range in {
			currentID, currentLocation, err := parseRow(row)

			if err != nil {
				outerror <- fmt.Errorf("Corrupted data passed to parser; error: %v", err)
			}

			if len(locationsInRide) != 0 && id != currentID {
				out <- model.Ride{ID: id, LocationSignals: locationsInRide}
				locationsInRide = []model.LocationSignal{}
			}

			id = currentID
			locationsInRide = append(locationsInRide, *currentLocation)
		}
		out <- model.Ride{ID: id, LocationSignals: locationsInRide} // handles the final row
		close(out)
	}()
	return out, outerror
}

func parseRow(row []string) (string, *model.LocationSignal, error) {
	if len(row) < 4 {
		return "", nil, fmt.Errorf("Row not containing 4 elemetns: %v", row)
	}

	id := row[0]
	latitude, errLat := strconv.ParseFloat(row[1], 64)
	longitude, errLon := strconv.ParseFloat(row[2], 64)
	timestamp, errTime := strconv.ParseInt(row[3], 10, 32)

	if errLat != nil || errLon != nil || errTime != nil {
		return "", nil, fmt.Errorf("Failed to parse row")
	}

	return id, &model.LocationSignal{
		Latitude:  latitude,
		Longitude: longitude,
		Timestamp: int32(timestamp),
	}, nil
}
