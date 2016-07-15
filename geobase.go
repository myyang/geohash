package geohash

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
	GetCenter() (float64, float64)
	ErrPair() (latErr, lngErr float64)
	Geohash() string
}

// LocationBox stores rectangle shape lcoation info and supports bounding box
type LocationBox struct {
	MaxLat, MinLat, MaxLng, MinLng float64
	LatErr, LngErr                 float64
	Hash                           string
	Precision                      int
}

// GetCenter return center of rectangle
func (lb LocationBox) GetCenter() (float64, float64) {
	return roundFloat64((lb.MaxLat+lb.MinLat)/2, lb.Precision), roundFloat64((lb.MaxLng+lb.MinLng)/2, lb.Precision)
}

// ErrPair returns estimated error distance pair in degree
func (lb LocationBox) ErrPair() (float64, float64) {
	return lb.LatErr, lb.LngErr
}

// Geohash return hash value
func (lb LocationBox) Geohash() string {
	return lb.Hash
}
