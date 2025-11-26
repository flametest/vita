package vo

import (
	"encoding/json"
	"math"

	"github.com/mmcloughlin/geohash"
	"github.com/shopspring/decimal"
)

// PI mathematical constant for calculations
const PI float64 = 3.141592653589793

// EarthRadius radius of Earth in meters (6378.137 km converted to meters)
const EarthRadius float64 = 6378.137 * 1000

// Coordinate represents a geographic coordinate with latitude and longitude
// using decimal.Decimal for precise numeric representation
type Coordinate struct {
	latitude  decimal.Decimal // geographic latitude
	longitude decimal.Decimal // geographic longitude
}

// NewCoordinate creates a new Coordinate instance with given latitude and longitude
func NewCoordinate(latitude decimal.Decimal, longitude decimal.Decimal) *Coordinate {
	return &Coordinate{latitude: latitude, longitude: longitude}
}

// Latitude returns the latitude of the coordinate
func (c *Coordinate) Latitude() decimal.Decimal {
	return c.latitude
}

// Longitude returns the longitude of the coordinate
func (c *Coordinate) Longitude() decimal.Decimal {
	return c.longitude
}

// MarshalJSON implements the json.Marshaler interface for Coordinate
func (c *Coordinate) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Latitude  decimal.Decimal `json:"latitude"`
		Longitude decimal.Decimal `json:"longitude"`
	}{
		Latitude:  c.latitude,
		Longitude: c.longitude,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface for Coordinate
func (c *Coordinate) UnmarshalJSON(bytes []byte) error {
	var tmp struct {
		Latitude  decimal.Decimal `json:"latitude"`
		Longitude decimal.Decimal `json:"longitude"`
	}
	err := json.Unmarshal(bytes, &tmp)
	if err != nil {
		return err
	}
	c.latitude = tmp.Latitude
	c.longitude = tmp.Longitude
	return nil
}

// NewCoordinateFromString creates a new Coordinate from latitude and longitude strings
func NewCoordinateFromString(latitude string, longitude string) (*Coordinate, error) {
	latitude1, err := decimal.NewFromString(latitude)
	if err != nil {
		return nil, err
	}
	longitude1, err := decimal.NewFromString(longitude)
	if err != nil {
		return nil, err
	}
	return NewCoordinate(latitude1, longitude1), nil
}

// NewCoordinateFromFloat64 creates a new Coordinate from latitude and longitude float64 values
func NewCoordinateFromFloat64(lat, lng float64) *Coordinate {
	lat1 := decimal.NewFromFloat(lat)
	lng1 := decimal.NewFromFloat(lng)
	return NewCoordinate(lat1, lng1)
}

// NewCoordinateFromGeoHash creates a new Coordinate from a geohash string
func NewCoordinateFromGeoHash(geoHash string) *Coordinate {
	lat, lng := geohash.Decode(geoHash)
	return NewCoordinateFromFloat64(lat, lng)
}

// DistanceFrom calculates the distance between this coordinate and another coordinate
// using the Haversine formula, returns distance in meters
func (c *Coordinate) DistanceFrom(vo *Coordinate) float64 {
	lat1, _ := c.latitude.Mul(decimal.NewFromFloat(PI)).Div(decimal.NewFromInt(180)).Float64()
	lat2, _ := vo.latitude.Mul(decimal.NewFromFloat(PI)).Div(decimal.NewFromInt(180)).Float64()
	lng1, _ := c.longitude.Mul(decimal.NewFromFloat(PI)).Div(decimal.NewFromInt(180)).Float64()
	lng2, _ := vo.longitude.Mul(decimal.NewFromFloat(PI)).Div(decimal.NewFromInt(180)).Float64()
	latDiff := lat1 - lat2
	lngDiff := lng1 - lng2
	return 2 * math.Asin(
		math.Sqrt(
			math.Pow(math.Sin(latDiff/2), 2)+math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(lngDiff/2), 2),
		),
	) * EarthRadius
}

// CalGeoHash calculates and returns the geohash string for this coordinate
func (c *Coordinate) CalGeoHash() string {
	latitude, _ := c.latitude.Float64()
	longitude, _ := c.longitude.Float64()
	return geohash.Encode(latitude, longitude)
}
