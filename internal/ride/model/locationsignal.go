package model

// LocationSignal provides information by a single reported location (latitude and logintude) at specific time (unix timestamp)
type LocationSignal struct {
	Latitude  float64
	Longitude float64
	Timestamp int32
}
