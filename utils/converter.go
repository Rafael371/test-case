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

	NormalizedDivisor = SecondsPerDayOnRoketin / SecondsPerDayOnEarth
)

func EarthToTotalSeconds(h, m, s int) int {
	return h*3600 + m*60 + s
}

func TotalSecondsToRoketin(total int) (int, int, int) {
	roketinTotalSeconds := math.Ceil(float64(total) * NormalizedDivisor)

	h := int(roketinTotalSeconds) / (MinutesOnRoketin * SecondOnRoketin)
	m := (int(roketinTotalSeconds) / SecondOnRoketin) % MinutesOnRoketin
	s := int(roketinTotalSeconds) % SecondOnRoketin

	return h, m, s
}