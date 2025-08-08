package utils

import "math"

const (
	HourOnRoketin    = 10
	MinutesOnRoketin = 100
	SecondOnRoketin  = 100

	HourOnEarth    = 24
	MinutesOnEarth = 60
	SecondOnEarth  = 60
)

var (
	SecondsPerDayOnEarth   = float64(HourOnEarth) * float64(MinutesOnEarth) * float64(SecondOnEarth)
	SecondsPerDayOnRoketin = float64(HourOnRoketin) * float64(MinutesOnRoketin) * float64(SecondOnRoketin)

	//ratio earth to Roketin hour
	NormalizedDivisor = SecondsPerDayOnRoketin / SecondsPerDayOnEarth
)

//Get earth total seconds
func EarthToTotalSeconds(h, m, s int) int {
	return h*3600 + m*60 + s
}

//Convert to Roketin hour
func ConvertToRoketinHour(total int) (int, int, int) {
	//Make total second on earth become total second on Roketin
	roketinTotalSeconds := math.Ceil(float64(total) * NormalizedDivisor)

	//Convert total second to hour, minute, and second
	h := int(roketinTotalSeconds) / (MinutesOnRoketin * SecondOnRoketin)
	m := (int(roketinTotalSeconds) / SecondOnRoketin) % MinutesOnRoketin
	s := int(roketinTotalSeconds) % SecondOnRoketin

	return h, m, s
}