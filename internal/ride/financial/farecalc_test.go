package financial

import (
	"testing"

	"github.com/kanzafirov/fareestimation/internal/ride/model"
)

type CalculateDistanceHaversineResult struct {
	lonFrom  float64
	latFrom  float64
	lonTo    float64
	latTo    float64
	expected float64
}

var distanceResults = []CalculateDistanceHaversineResult{
	{37.96666, 23.728308, 37.966627, 23.728263, 0.006026788669283551},
	{37.96666, 23.728308, 37.96666, 23.728308, 0},
	{38.017711, 23.834016, 38.001076, 23.787746, 5.416155515757538},
	{37.969647, 23.722431, 38.001076, 23.787746, 7.935877881075722},
	{37.969647, 23.722431, 43.172348, 27.926977, 699.646090855085},
}

func TestCalculateDistanceHaversine(t *testing.T) {
	for _, test := range distanceResults {
		result := calculateDistanceHaversine(test.lonFrom, test.latFrom, test.lonTo, test.latTo)

		if result != test.expected {
			t.Errorf("[TestCalculateDistanceHaversine] Failed: expected: %v, actual: %v", test.expected, result)
		}
	}
}

type IsDayShiftResult struct {
	timestamp int32
	expected  bool
}

var isDayShiftResult = []IsDayShiftResult{
	{1588439749, true},
	{1405030266, true},
	{1538441931, false},
	{1388551388, false},
}

func TestIsDayShift(t *testing.T) {
	for _, test := range isDayShiftResult {
		result := isDayShift(test.timestamp)

		if result != test.expected {
			t.Errorf("[TestIsDayShift] Failed: expected: %t, actual: %t, timestamp: %d", test.expected, result, test.timestamp)
		}
	}
}

type CalculateTotalFareEstimateResult struct {
	ride     *model.Ride
	expected *FareEstimate
}

var fareEstimateResults = []CalculateTotalFareEstimateResult{
	{&model.Ride{ID: "a", LocationSignals: []model.LocationSignal{{37.966660, 23.728308, 1405594957}, {37.935490, 23.625655, 1405596220}}}, &FareEstimate{IDRide: "a", Total: 8.43730278512446}},
	{&model.Ride{ID: "1234", LocationSignals: []model.LocationSignal{{37.966660, 23.728308, 1405594957}, {37.935490, 23.625655, 1405596220}}}, &FareEstimate{IDRide: "1234", Total: 8.43730278512446}},
	{&model.Ride{ID: "1234", LocationSignals: []model.LocationSignal{{37.966660, 23.728308, 1405594957}, {37.935490, 23.625655, 1405596220}, {37.966195, 23.728613, 1405595043}, {37.964823, 23.726953, 1405595102}}}, &FareEstimate{IDRide: "1234", Total: 4.702681277042642}},
	{&model.Ride{ID: "6", LocationSignals: []model.LocationSignal{{37.966660, 23.728308, 1405594957}, {37.954302, 23.713370, 1405595284}, {37.938042, 23.692308, 1405595362}, {37.938985, 23.690435, 1405595371}, {37.964823, 23.726953, 1405595102}}}, &FareEstimate{IDRide: "6", Total: 3.47}}, // corrupted point + min fee
}

func TestCalculateTotalFareEstimate(t *testing.T) {
	for _, test := range fareEstimateResults {
		result := calculateTotalFareEstimate(test.ride)

		if *result != *test.expected {
			t.Errorf("[TestCalculateTotalFareEstimate] Failed: expected: %v, actual: %v", *test.expected, *result)
		}
	}
}
