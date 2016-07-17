package geohash

import "bytes"
import "fmt"

// fix constants
const (
	initUnitLat   = 30.0
	initUnitLng   = 60.0
	DefaultB36Str = "23456789bBCdDFgGhHjJKlLMnNPqQrRtTVWX"
)

// Helper variables
var (
	B36 = []byte(DefaultB36Str)
)

func encode36(latitude, longitude float64, key []byte, precision int) (
	hash string, maxLat, minLat, maxLng, minLng, unitLat, unitLng float64) {
	b := make([]byte, 0, precision)
	unitLng, maxLng, minLng = initUnitLng, MaxLng, MinLng
	unitLat, maxLat, minLat = initUnitLat, MaxLat, MinLat
	for i := 0; i < precision; i++ {
		// since max search time is 12 (6 for each),
		// just iterate it instead of binary search
		clat, clng := 0, 0
		for (maxLat - unitLat) > latitude {
			clat++
			maxLat -= unitLat
			// float number present limitation &&
			// we divide unit by 6 which contains 3 leads to repeating decimal
			if clat == 6 && maxLat-latitude < 10e-9 {
				clat--
			}
		}
		for (minLng + unitLng) < longitude {
			clng++
			minLng += unitLng
			if clng == 6 && maxLat-latitude < 10e-9 {
				clat--
			}
		}
		// update square
		minLat, maxLng = maxLat-unitLat, minLng+unitLng
		unitLat, unitLng = unitLat/6, unitLng/6
		b = append(b, key[clat*6+clng])
	}
	_, _ = (minLat+maxLat)/2, (minLng+maxLng)/2
	return string(b[:]), maxLat, minLat, maxLng, minLng, roundFloat64(unitLat, precision), roundFloat64(unitLng, precision)
}

func decode36(hashv string, key []byte, precision int) (
	lat, lng, maxLat, minLat, maxLng, minLng, uLat, uLng float64) {
	hashb := []byte(hashv)
	if precision <= 0 {
		precision = len(hashb)
	}
	uLat, maxLat, minLat = initUnitLat, MaxLat, MinLat
	uLng, maxLng, minLng = initUnitLng, MaxLng, MinLng
	for _, v := range hashb {
		i := bytes.IndexByte(B36, v)
		row, col := i/6, i%6
		maxLat = maxLat - float64(row)*uLat
		minLat = maxLat - uLat
		minLng = minLng + float64(col)*uLng
		maxLng = minLng + uLng
		uLat, uLng = uLat/6, uLng/6
	}
	return roundFloat64((maxLat+minLat)/2, precision), roundFloat64((maxLng+minLng)/2, precision),
		maxLat, minLat, maxLng, minLng, uLat, uLng
}

// NewDefaultGeoHash36 generate geohash36 cryptor with default key
func NewDefaultGeoHash36() GeoCryptor {
	return NewGeoHash36(DefaultB36Str)
}

// NewGeoHash36 generate geohash36 with given key
func NewGeoHash36(key string) GeoCryptor {
	g := &GeoHash36{}
	g.SetKey(key)
	return g
}

// GeoHash36 is one of GeoCryptor that provide geohash36 coding
// for more detai, please check wiki page:
// https://en.wikipedia.org/wiki/Geohash-36
// or promotion website: http://geo36.org
type GeoHash36 struct {
	key []byte
}

// SetKey set hash key value
func (g *GeoHash36) SetKey(key string) {
	g.key = []byte(key)
}

// HashKey return hash key of this hasher
func (g *GeoHash36) HashKey() string {
	return string(g.key)
}

// Encode and return hash value only
func (g *GeoHash36) Encode(latitude, longitude float64, precision int) string {
	v, _, _, _, _, _, _ := encode36(latitude, longitude, g.key, precision)
	return v
}

// Decode and return central lat, lng pair
func (g *GeoHash36) Decode(value string, precision int) (float64, float64) {
	lat, lng, _, _, _, _, _, _ := decode36(value, g.key, precision)
	return lat, lng
}

// EncodeWithErr returns also estimate error in degree
func (g *GeoHash36) EncodeWithErr(latitude, longitude float64, precision int) (string, float64, float64) {
	v, _, _, _, _, latErr, lngErr := encode36(latitude, longitude, g.key, precision)
	return v, latErr * 6, lngErr * 6
}

// DecodeWithErr returns also estimate error
func (g *GeoHash36) DecodeWithErr(value string, precision int) (float64, float64, float64, float64) {
	lat, lng, _, _, _, _, latErr, lngErr := decode36(value, g.key, precision)
	return lat, lng, latErr * 6, lngErr * 6
}

// EncodeAsBox returns a location box
func (g *GeoHash36) EncodeAsBox(latitude, longitude float64, precision int) BoundingBox {
	v, maxlat, minlat, maxlng, minlng, latErr, lngErr := encode36(latitude, longitude, g.key, precision)
	return &LocationBox{
		MaxLat: maxlat, MinLat: minlat, MaxLng: maxlng, MinLng: minlng,
		LatErr: latErr * 6, LngErr: lngErr * 6, Hash: v, Precision: precision}
}

// DecodeAsBox returns a location box
func (g *GeoHash36) DecodeAsBox(value string, precision int) BoundingBox {
	_, _, maxlat, minlat, maxlng, minlng, latErr, lngErr := decode36(value, g.key, precision)
	return &LocationBox{
		MaxLat: maxlat, MinLat: minlat, MaxLng: maxlng, MinLng: minlng,
		LatErr: latErr * 6, LngErr: lngErr * 6, Hash: value, Precision: precision}
}

// Neighbors returns adjcent 8 neighbors
func (g *GeoHash36) Neighbors(value string, precision int) []BoundingBox {
	lb := g.DecodeAsBox(value, precision)
	neighbor := func(lb *LocationBox, dlat, dlng int) *LocationBox {
		latUnit, lngUnit := (lb.MaxLat-lb.MinLat)*float64(dlat), (lb.MaxLng-lb.MinLng)*float64(dlng)
		r := &LocationBox{
			MaxLat: lb.MaxLat + latUnit, MinLat: lb.MinLat + latUnit,
			MaxLng: lb.MaxLng + lngUnit, MinLng: lb.MinLng + lngUnit,
			LatErr: lb.LatErr, LngErr: lb.LngErr}
		lat, lng, err := r.GetCenter()
		if err != nil {
			fmt.Printf("Neighbor error: %v\n", err)
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
