package geohash

import (
	"fmt"
)

// fixed constants
const (
	MaxLat float64 = 90.0
	MinLat         = -90.0
	MaxLng         = 180.0
	MinLng         = -180.0
)

// GeoCryptor defines geohash provider
type GeoCryptor interface {
	HashKey() string
	SetKey(hashkey string)
	Encode(latitude, longitude float64, precision int) string
	EncodeWithErr(latitude, longitude float64, precision int) (hash string, latErr, lngErr float64)
	EncodeAsBox(latitude, longitude float64, precision int) BoundingBox
	Decode(value string, precision int) (lat, lng float64)
	DecodeWithErr(value string, precision int) (lat, lng, latErr, lngErr float64)
	DecodeAsBox(value string, precision int) BoundingBox
	Neighbors(value string, precision int) []BoundingBox
}

// BoundingBox caculate shape center
type BoundingBox interface {
	GetCenter() (float64, float64, error)
	ErrPair() (latErr, lngErr float64)
	Geohash() (string, error)
}

// LocationBox stores rectangle shape lcoation info and supports bounding box
type LocationBox struct {
	MaxLat, MinLat, MaxLng, MinLng float64
	LatErr, LngErr                 float64
	Hash                           string
	Precision                      int
}

// GetCenter return center of rectangle
func (lb LocationBox) GetCenter() (float64, float64, error) {
	if err := lb.validBox(); err != nil {
		return -999.0, -999.0, err
	}
	return roundFloat64((lb.MaxLat+lb.MinLat)/2, lb.Precision), roundFloat64((lb.MaxLng+lb.MinLng)/2, lb.Precision), nil
}

// ErrPair returns estimated error distance pair in degree
func (lb LocationBox) ErrPair() (float64, float64) {
	return lb.LatErr, lb.LngErr
}

// Geohash return hash value
func (lb LocationBox) Geohash() (string, error) {
	if err := lb.validBox(); err != nil {
		return "----", err
	}
	return lb.Hash, nil
}

func (lb LocationBox) validBox() error {
	if lb.MaxLat > MaxLat || lb.MinLat > MaxLat || lb.MinLat < MinLat || lb.MaxLat < MinLat {
		return CoordinateError{LB: lb}
	}

	if lb.MaxLng > MaxLng || lb.MinLng > MaxLng {
		if lb.MaxLng > lb.MinLng {
			lb.MinLng = lb.MinLng + MinLng*2
			lb.MaxLng = lb.MinLng + lb.LngErr*2
		} else {
			return CoordinateError{LB: lb}
		}
	} else if lb.MinLng < MinLng || lb.MaxLng < MinLng {
		if lb.MinLng < lb.MaxLng {
			lb.MaxLng = lb.MaxLng + MaxLng*2
			lb.MinLng = lb.MinLng - lb.LngErr*2
		} else {
			return CoordinateError{LB: lb}
		}
	}
	return nil
}

// CoordinateError isa custom error
type CoordinateError struct {
	LB LocationBox
}

func (ce CoordinateError) Error() string {
	return fmt.Sprintf(
		"Wrong coordinate:\nmaxLat: %v,  minLat: %v, maxLng: %v, minLng: %v",
		ce.LB.MaxLat, ce.LB.MinLat, ce.LB.MaxLng, ce.LB.MinLng)
}
