package geohash

import "bytes"

const (
	initUnitLat = 30.0
	initUnitLng = 60.0
)

// Helper variables
var (
	B36 = []byte("23456789bBCdDFgGhHjJKlLMnNPqQrRtTVWX")
)

// Encode36 provides geohash-36 encode
// return the encoded string, latErr, lngErr in degree
func Encode36(latitude, longitude float64, precision int) (string, float64, float64) {
	b := make([]byte, 0, precision)
	unitLng, maxLng, minLng := initUnitLng, MaxLng, MinLng
	unitLat, maxLat, minLat := initUnitLat, MaxLat, MinLat
	for i := 0; i < precision; i++ {
		// since max search time is 12 (6 for each),
		// just iterate it instead of binary search
		clat, clng := 0, 0
		for (maxLat - unitLat) > latitude {
			clat++
			maxLat -= unitLat
		}
		for (minLng + unitLng) < longitude {
			clng++
			minLng += unitLng
		}
		// update square
		minLat, maxLng = maxLat-unitLat, minLng+unitLng
		unitLat, unitLng = unitLat/6.0, unitLng/6.0
		b = append(b, B36[clat*6+clng])
	}
	_, _ = (minLat+maxLat)/2, (minLng+maxLng)/2
	return string(b[:]), unitLat, unitLng
}

// Decode36 decode given hash value into latitude and longitude centry point
func Decode36(hashv string, precision int) (float64, float64) {
	hashb := []byte(hashv)
	if precision <= 0 {
		precision = len(hashb)
	}
	uLat, maxLat, minLat := initUnitLat, MaxLat, MinLat
	uLng, maxLng, minLng := initUnitLng, MaxLng, MinLng
	for _, v := range hashb {
		i := bytes.IndexByte(B36, v)
		row, col := i/6, i%6
		maxLat = maxLat - float64(row)*uLat
		minLat = maxLat - uLat
		minLng = minLng + float64(col)*uLng
		maxLng = minLng + uLng
		uLat, uLng = uLat/6, uLng/6
	}
	return roundFloat64((maxLat+minLat)/2, precision), roundFloat64((maxLng+minLng)/2, precision)
}
