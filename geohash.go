package geohash

import "bytes"

import "log"

// fixed constants
const (
	ByteWidth     int = 4
	DefaultB32Str     = "0123456789bcdefghjkmnpqrstuvwxyz"
)

// Helper variables
var (
	B32           = []byte(DefaultB32Str)
	Bits          = []int{16, 8, 4, 2, 1, 0}
	HalfByteWidth = (ByteWidth / 2) + 1
)

func encode(latitude, longitude float64,
	key []byte, precision int) (string, float64, float64, float64, float64) {

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
			bf = append(bf, key[ch])
			ch, byteCount, resultLen = 0, 0, resultLen+1
		}
	}

	return string(bf[:]), maxLat, minLat, maxLng, minLng
}

func decode(hashv string, key []byte, precision int) (float64, float64, float64, float64, float64, float64) {
	hashb := []byte(hashv)
	if precision <= 0 {
		precision = len(hashb)
	}
	minLat, maxLat, minLng, maxLng := MinLat, MaxLat, MinLng, MaxLng
	alter := true
	for i := 0; i < len(hashb); i++ {
		byteCount := 0
		v := bytes.IndexByte(key, hashb[i])
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
	return roundFloat64(lat, precision), roundFloat64(lng, precision), maxLat, minLat, maxLng, minLng
}

func latErr(l int) float64 {
	b := 2 << uint(2*l+l/2)
	return roundFloat64(180.0/float64(b), l)
}

func lngErr(l int) float64 {
	b := 2 << uint(3*l-l/2)
	return roundFloat64(360.0/float64(b), l)
}

// NewDefaultGeoHash return a geohash cryptor with defined key
func NewDefaultGeoHash() GeoCryptor {
	return NewGeoHash(DefaultB32Str)
}

// NewGeoHash return a geohash cryptor with given key
func NewGeoHash(key string) GeoCryptor {
	g := &GeoHash{}
	g.SetKey(key)
	return g
}

// GeoHash provides functions to compute geohash value
// for more detail, please check following link
// https://en.wikipedia.org/wiki/Geohash
type GeoHash struct {
	key []byte
}

// SetKey set hash key value
func (g *GeoHash) SetKey(key string) {
	g.key = []byte(key)
}

// HashKey return hash key of this hasher
func (g *GeoHash) HashKey() string {
	return string(g.key)
}

// Encode and return hash value only
func (g *GeoHash) Encode(latitude, longitude float64, precision int) string {
	v, _, _, _, _ := encode(latitude, longitude, g.key, precision)
	return v
}

// Decode and return central lat, lng pair
func (g *GeoHash) Decode(value string, precision int) (float64, float64) {
	lat, lng, _, _, _, _ := decode(value, g.key, precision)
	return lat, lng
}

// EncodeWithErr returns also estimate error in degree
func (g *GeoHash) EncodeWithErr(latitude, longitude float64, precision int) (string, float64, float64) {
	v, _, _, _, _ := encode(latitude, longitude, g.key, precision)
	return v, latErr(precision), lngErr(precision)
}

// DecodeWithErr returns also estimate error
func (g *GeoHash) DecodeWithErr(value string, precision int) (float64, float64, float64, float64) {
	lat, lng, _, _, _, _ := decode(value, g.key, precision)
	return lat, lng, latErr(precision), lngErr(precision)
}

// EncodeAsBox returns a location box
func (g *GeoHash) EncodeAsBox(latitude, longitude float64, precision int) BoundingBox {
	v, maxlat, minlat, maxlng, minlng := encode(latitude, longitude, g.key, precision)
	return &LocationBox{
		MaxLat: maxlat, MinLat: minlat, MaxLng: maxlng, MinLng: minlng,
		LatErr: latErr(precision), LngErr: lngErr(precision), Hash: v, Precision: precision}
}

// DecodeAsBox returns a location box
func (g *GeoHash) DecodeAsBox(value string, precision int) BoundingBox {
	_, _, maxlat, minlat, maxlng, minlng := decode(value, g.key, precision)
	return &LocationBox{
		MaxLat: maxlat, MinLat: minlat, MaxLng: maxlng, MinLng: minlng,
		LatErr: latErr(precision), LngErr: lngErr(precision), Hash: value, Precision: precision}
}

// Neighbors returns adjcent 8 neighbors
func (g *GeoHash) Neighbors(value string, precision int) []BoundingBox {
	lb := g.DecodeAsBox(value, precision)
	neighbor := func(lb *LocationBox, dlat, dlng int) *LocationBox {
		latUnit, lngUnit := lb.LatErr*2*float64(dlat), lb.LngErr*2*float64(dlng)
		r := &LocationBox{
			MaxLat: lb.MaxLat + latUnit, MinLat: lb.MinLat + latUnit,
			MaxLng: lb.MaxLng + lngUnit, MinLng: lb.MinLng + lngUnit,
			LatErr: lb.LatErr, LngErr: lb.LngErr, Precision: lb.Precision}
		lat, lng, err := r.GetCenter()
		if err != nil {
			log.Printf("Erro while get neighbor <%v>\n", err)
			return nil
		}
		r.Hash = g.Encode(lat, lng, precision)
		return r
	}
	n := make([]BoundingBox, 0, 8)
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if i == 0 && j == 0 {
				continue
			}
			n = append(n, neighbor(lb.(*LocationBox), i, j))
		}
	}
	return n
}
