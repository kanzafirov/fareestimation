package model

// Ride is a structure providing all the location signals repored by a driver with specific ride id
type Ride struct {
	ID              string
	LocationSignals []LocationSignal
}
