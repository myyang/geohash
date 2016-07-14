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
	var b bytes.Buffer
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
		b.WriteByte(B36[clat*6+clng])
	}
	_, _ = (minLat+maxLat)/2, (minLng+maxLng)/2
	return b.String(), unitLat, unitLng
}

// Decode36 decode given hash value into latitude and longitude centry point
func Decode36(hashv string, precision int) (float64, float64) {
	hashb := []byte(hashv)
	if precision <= 0 {
		precision = len(hashb)
	}
	m := make([]int, len(hashb))
	for i, v := range hashb {
		m[i] = bytes.IndexByte(B36, v)
	}
	lat, _, _ := decodeLat(m)
	lng, _, _ := decodeLng(m)
	return roundFloat64(lat, precision), roundFloat64(lng, precision)
}

func decodeLat(m []int) (float64, float64, float64) {
	unit, max, min := initUnitLat, MaxLat, MinLat
	for _, v := range m {
		offset := v / 6
		max = max - float64(offset)*unit
		min = max - unit
		unit /= 6.0
	}

	return (max + min) / 2, max, min
}

func decodeLng(m []int) (float64, float64, float64) {
	unit, max, min := initUnitLng, MaxLng, MinLng
	for _, v := range m {
		offset := v % 6
		min = min + float64(offset)*unit
		max = min + unit
		unit /= 6.0
	}
	return (max + min) / 2, max, min
}
