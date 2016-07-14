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

	var bf bytes.Buffer
	latChan := l2bChan(latitude, MaxLat, MinLat, precision*HalfByteWidth)
	lngChan := l2bChan(longitude, MaxLng, MinLng, precision*HalfByteWidth)

	resultLen, ch, byteCount, alter := 0, 0, 0, true
	for resultLen < precision {
		switch alter {
		case true:
			if v := <-lngChan; v == 1 {
				ch |= Bits[byteCount]
			}
		case false:
			if v := <-latChan; v == 1 {
				ch |= Bits[byteCount]
			}
		}
		alter = !alter
		if byteCount < ByteWidth {
			byteCount++
		} else {
			bf.WriteByte(B32[ch])
			ch, byteCount, resultLen = 0, 0, resultLen+1
		}
	}

	return bf.String()
}

func l2bChan(loc, max, min float64, subPrecision int) <-chan int {
	rchan := make(chan int, subPrecision)
	go func() {
		for subPrecision > 0 {
			if mid := (max + min) / 2.0; mid < loc {
				rchan <- 1
				min = mid
			} else {
				rchan <- 0
				max = mid
			}
			subPrecision--
		}
	}()
	return rchan
}

// Decode given string and return (lat, lng) pair
func Decode(hashv string, precision int) (float64, float64) {
	hashb := []byte(hashv)
	if precision <= 0 {
		precision = len(hashb)
	}
	toLatChan, toLngChan := make(chan int, precision), make(chan int, precision)
	latChan, lngChan := b2lChan(toLatChan, MaxLat, MinLat), b2lChan(toLngChan, MaxLng, MinLng)
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
				toLngChan <- b
			case false:
				toLatChan <- b
			}
			alter = !alter
			byteCount++
		}
	}
	close(toLngChan)
	close(toLatChan)
	lat, lng := <-latChan, <-lngChan
	return roundFloat64(lat, precision), roundFloat64(lng, precision)
}

func b2lChan(bChan chan int, max, min float64) <-chan float64 {
	rchan := make(chan float64)
	go func() {
		mid := 0.0
		for {
			b, ok := <-bChan
			if !ok {
				break
			}
			mid = (max + min) / 2
			if b == 1 {
				min = mid
			} else {
				max = mid
			}
		}
		mid = (max + min) / 2
		rchan <- mid
	}()
	return rchan
}
