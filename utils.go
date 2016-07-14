package geohash

import "math"

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func roundFloat64(num float64, precision int) float64 {
	base := math.Pow(10, float64(precision))
	return float64(round(num*base)) / base
}
