package financial

import (
	"fmt"
	"math"
	"time"

	"github.com/kanzafirov/fareestimation/internal/ride/model"
)

const buffSize = 50

const (
	earthRadius          = float64(6371)
	minSpeed             = float64(10)
	maxSpeed             = float64(100)
	flagRate             = float64(1.3)
	idleHourRate         = float64(11.9)
	movingRateDayShift   = float64(0.74)
	movingRateNightShift = float64(1.3)
	minTotal             = float64(3.47)
)

// NewFairFeeCalculationChannel gets an input channel as parameter, passing Ride objects, for each of them
// a fare rate is calculated and returns a new []string channel, where passes id and estimation fee
func NewFairFeeCalculationChannel(in <-chan model.Ride) <-chan []string {
	out := make(chan []string, buffSize)

	go func() {
		for ride := range in {
			estimationFee := calculateTotalFareEstimate(&ride)
			out <- []string{estimationFee.IDRide, fmt.Sprintf("%.2f", estimationFee.Total)}
		}
		close(out)
	}()
	return out
}

func calculateTotalFareEstimate(ride *model.Ride) *FareEstimate {
	total := flagRate

	for i := 0; i < len(ride.LocationSignals)-1; i++ {
		for j := i + 1; j < len(ride.LocationSignals); j++ {

			startSignal := ride.LocationSignals[i]
			endSignal := ride.LocationSignals[i+1]

			deltaTimeSeconds := float64(endSignal.Timestamp - startSignal.Timestamp)
			deltaDistanceKm := calculateDistanceHaversine(startSignal.Longitude, startSignal.Latitude, endSignal.Longitude, endSignal.Latitude)

			speedKmh := (deltaDistanceKm / deltaTimeSeconds) * 3600

			if speedKmh > maxSpeed {
				// set the next start signal to one next so it skips the corrupted
				i++
				// skip the corrupted point
				continue
			}

			if speedKmh <= minSpeed {
				// calculate idle rate
				total += (deltaTimeSeconds / 3600) * idleHourRate
				break
			}

			// calculate distance rate by hour
			if isDayShift(startSignal.Timestamp) {
				total += deltaDistanceKm * movingRateDayShift
			} else {
				total += deltaDistanceKm * movingRateNightShift
			}
			break
		}
	}

	total = math.Max(total, minTotal)

	return &FareEstimate{
		IDRide: ride.ID,
		Total:  total,
	}
}

func isDayShift(timestamp int32) bool {
	t := time.Unix(int64(timestamp), 0).UTC()
	hour := t.Hour()

	if hour >= 5 && hour < 24 {
		return true
	}

	return false
}

// calculateDistanceHaversine will calculate the spherical distance as the
// crow flies between lat and lon for two given points in km by the Haverstine formula
func calculateDistanceHaversine(lonFrom float64, latFrom float64, lonTo float64, latTo float64) (distance float64) {
	var deltaLat = (latTo - latFrom) * (math.Pi / 180)
	var deltaLon = (lonTo - lonFrom) * (math.Pi / 180)

	var a = math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(latFrom*(math.Pi/180))*math.Cos(latTo*(math.Pi/180))*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance = earthRadius * c

	return
}
