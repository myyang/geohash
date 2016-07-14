package geohash

import "bytes"

// fixed constants
const (
	ByteWidth int = 4
)

// Helper variables
var (
	B32           = []byte("0123456789bcdefghjkmnpqrstuvwxyz")
	Bits          = []int{16, 8, 4, 2, 1, 0}
	HalfByteWidth = (ByteWidth / 2) + 1
)

// Encode given latitude and longitude and return a encoded string with
// given precision (aka encoded string length)
func Encode(latitude, longitude float64, precision int) string {
	if latitude > MaxLat || latitude < MinLat || longitude > MaxLng || longitude < MinLng || precision > 12 {
		return ""
	}

	minLat, maxLat, minLng, maxLng := MinLat, MaxLat, MinLng, MaxLng
	bf := make([]byte, 0, precision)

	resultLen, ch, byteCount, alter := 0, 0, 0, true
	for resultLen < precision {
		switch alter {
		case true:
			if midLng := (minLng + maxLng) / 2; midLng < longitude {
				ch |= Bits[byteCount]
				minLng = midLng
			} else {
				maxLng = midLng
			}
		case false:
			if midLat := (minLat + maxLat) / 2; midLat < latitude {
				ch |= Bits[byteCount]
				minLat = midLat
			} else {
				maxLat = midLat
			}
		}
		alter = !alter
		if byteCount < ByteWidth {
			byteCount++
		} else {
			bf = append(bf, B32[ch])
			ch, byteCount, resultLen = 0, 0, resultLen+1
		}
	}

	return string(bf[:])
}

// Decode given string and return (lat, lng) pair
func Decode(hashv string, precision int) (float64, float64) {
	hashb := []byte(hashv)
	if precision <= 0 {
		precision = len(hashb)
	}
	minLat, maxLat, minLng, maxLng := MinLat, MaxLat, MinLng, MaxLng
	alter := true
	for i := 0; i < len(hashb); i++ {
		byteCount := 0
		v := bytes.IndexByte(B32, hashb[i])
		for byteCount <= ByteWidth {
			b := 0
			if v&Bits[byteCount] > 0 {
				b = 1
				v -= Bits[byteCount]
			}
			switch alter {
			case true:
				if b == 1 {
					minLng = (maxLng + minLng) / 2
				} else {
					maxLng = (maxLng + minLng) / 2
				}
			case false:
				if b == 1 {
					minLat = (maxLat + minLat) / 2
				} else {
					maxLat = (maxLat + minLat) / 2
				}
			}
			alter = !alter
			byteCount++
		}
	}
	lat, lng := (maxLat+minLat)/2, (maxLng+minLng)/2
	return roundFloat64(lat, precision), roundFloat64(lng, precision)
}
