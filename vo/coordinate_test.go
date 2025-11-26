package vo

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"
)

func TestNewCoordinate(t *testing.T) {
	lat := decimal.NewFromFloat(40.7128)
	lng := decimal.NewFromFloat(-74.0060)

	coord := NewCoordinate(lat, lng)

	if coord.Latitude().Cmp(lat) != 0 {
		t.Errorf("Expected latitude %v, got %v", lat, coord.Latitude())
	}

	if coord.Longitude().Cmp(lng) != 0 {
		t.Errorf("Expected longitude %v, got %v", lng, coord.Longitude())
	}
}

func TestNewCoordinateFromString(t *testing.T) {
	tests := []struct {
		name      string
		latStr    string
		lngStr    string
		expectLat string
		expectLng string
		expectErr bool
	}{
		{
			name:      "Valid coordinates",
			latStr:    "40.7128",
			lngStr:    "-74.0060",
			expectLat: "40.7128",
			expectLng: "-74.006",
			expectErr: false,
		},
		{
			name:      "Invalid latitude",
			latStr:    "invalid",
			lngStr:    "-74.0060",
			expectErr: true,
		},
		{
			name:      "Invalid longitude",
			latStr:    "40.7128",
			lngStr:    "invalid",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coord, err := NewCoordinateFromString(tt.latStr, tt.lngStr)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if coord.Latitude().String() != tt.expectLat {
				t.Errorf("Expected latitude %s, got %s", tt.expectLat, coord.Latitude().String())
			}

			if coord.Longitude().String() != tt.expectLng {
				t.Errorf("Expected longitude %s, got %s", tt.expectLng, coord.Longitude().String())
			}
		})
	}
}

func TestNewCoordinateFromFloat64(t *testing.T) {
	lat := 40.7128
	lng := -74.0060

	coord := NewCoordinateFromFloat64(lat, lng)

	coordLat, _ := coord.Latitude().Float64()
	coordLng, _ := coord.Longitude().Float64()

	if coordLat != lat {
		t.Errorf("Expected latitude %v, got %v", lat, coordLat)
	}

	if coordLng != lng {
		t.Errorf("Expected longitude %v, got %v", lng, coordLng)
	}
}

func TestNewCoordinateFromGeoHash(t *testing.T) {
	// Using known geohash for New York City area
	geoHash := "dr5ruycj8z70"

	coord := NewCoordinateFromGeoHash(geoHash)

	// Validate that the coordinates are reasonable for NYC area
	lat, _ := coord.Latitude().Float64()
	lng, _ := coord.Longitude().Float64()

	if lat < 39 || lat > 41 {
		t.Errorf("Latitude %v is not in expected range for NYC area", lat)
	}

	if lng < -75 || lng > -73 {
		t.Errorf("Longitude %v is not in expected range for NYC area", lng)
	}

	// Test that geohash calculation returns the same value
	calculatedHash := coord.CalGeoHash()
	if calculatedHash[:7] != geoHash[:7] { // Compare first 7 characters due to precision
		t.Errorf("Expected geohash starting with %s, got %s", geoHash[:7], calculatedHash[:7])
	}
}

func TestCoordinateMarshalJSON(t *testing.T) {
	lat := decimal.NewFromFloat(40.7128)
	lng := decimal.NewFromFloat(-74.0060)
	coord := NewCoordinate(lat, lng)

	data, err := json.Marshal(coord)
	if err != nil {
		t.Fatalf("Failed to marshal coordinate: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["latitude"].(string) != lat.String() {
		t.Errorf("Expected latitude %s, got %v", lat.String(), result["latitude"])
	}

	if result["longitude"].(string) != lng.String() {
		t.Errorf("Expected longitude %s, got %v", lng.String(), result["longitude"])
	}
}

func TestCoordinateUnmarshalJSON(t *testing.T) {
	jsonData := `{"latitude": "40.7128", "longitude": "-74.0060"}`

	var coord Coordinate
	err := json.Unmarshal([]byte(jsonData), &coord)
	if err != nil {
		t.Fatalf("Failed to unmarshal coordinate: %v", err)
	}

	if coord.Latitude().String() != "40.7128" {
		t.Errorf("Expected latitude 40.7128, got %s", coord.Latitude().String())
	}

	if coord.Longitude().String() != "-74.006" {
		t.Errorf("Expected longitude -74.006, got %s", coord.Longitude().String())
	}
}

func TestCoordinateUnmarshalJSONInvalid(t *testing.T) {
	tests := []struct {
		name      string
		jsonData  string
		expectErr bool
	}{
		{
			name:      "Invalid JSON",
			jsonData:  `{"latitude": "invalid", "longitude": "-74.0060"}`,
			expectErr: true,
		},
		// Note: The json.Unmarshal will not fail for missing fields or empty JSON
		// because the fields are pointers and will have zero values
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var coord Coordinate
			err := json.Unmarshal([]byte(tt.jsonData), &coord)

			if tt.expectErr && err == nil {
				t.Errorf("Expected error, got nil")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestDistanceFrom(t *testing.T) {
	tests := []struct {
		name        string
		coord1Lat   float64
		coord1Lng   float64
		coord2Lat   float64
		coord2Lng   float64
		expectedMax float64 // Maximum expected distance in meters
	}{
		{
			name:        "Same coordinates",
			coord1Lat:   40.7128,
			coord1Lng:   -74.0060,
			coord2Lat:   40.7128,
			coord2Lng:   -74.0060,
			expectedMax: 1.0, // Very small distance
		},
		{
			name:        "New York to Los Angeles",
			coord1Lat:   40.7128,
			coord1Lng:   -74.0060,
			coord2Lat:   34.0522,
			coord2Lng:   -118.2437,
			expectedMax: 4000000.0, // Around 4000km
		},
		{
			name:        "Short distance (NYC to Brooklyn)",
			coord1Lat:   40.7128,
			coord1Lng:   -74.0060,
			coord2Lat:   40.6892,
			coord2Lng:   -73.9442,
			expectedMax: 10000.0, // Around 10km
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coord1 := NewCoordinateFromFloat64(tt.coord1Lat, tt.coord1Lng)
			coord2 := NewCoordinateFromFloat64(tt.coord2Lat, tt.coord2Lng)

			distance := coord1.DistanceFrom(coord2)

			if distance < 0 {
				t.Errorf("Distance should not be negative, got %f", distance)
			}

			if distance > tt.expectedMax {
				t.Errorf("Distance %f meters exceeds expected maximum %f meters", distance, tt.expectedMax)
			}
		})
	}
}

func TestCalGeoHash(t *testing.T) {
	tests := []struct {
		name     string
		lat      float64
		lng      float64
		expected string // First few characters of expected geohash
	}{
		{
			name:     "New York City",
			lat:      40.7128,
			lng:      -74.0060,
			expected: "dr5reg", // Updated expected geohash
		},
		{
			name:     "London",
			lat:      51.5074,
			lng:      -0.1278,
			expected: "gcpvj0", // Updated expected geohash
		},
		{
			name:     "Tokyo",
			lat:      35.6762,
			lng:      139.6503,
			expected: "xn76cy", // Updated expected geohash
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coord := NewCoordinateFromFloat64(tt.lat, tt.lng)
			geohash := coord.CalGeoHash()

			if len(geohash) < len(tt.expected) {
				t.Errorf("Geohash %s is shorter than expected", geohash)
			}

			if geohash[:len(tt.expected)] != tt.expected {
				t.Errorf("Expected geohash starting with %s, got %s", tt.expected, geohash[:len(tt.expected)])
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("Coordinates at equator and prime meridian", func(t *testing.T) {
		coord := NewCoordinateFromFloat64(0.0, 0.0)
		geohash := coord.CalGeoHash()

		if len(geohash) == 0 {
			t.Errorf("Geohash should not be empty for valid coordinates")
		}
	})

	t.Run("Coordinates at poles", func(t *testing.T) {
		northPole := NewCoordinateFromFloat64(90.0, 0.0)
		southPole := NewCoordinateFromFloat64(-90.0, 0.0)

		northHash := northPole.CalGeoHash()
		southHash := southPole.CalGeoHash()

		if len(northHash) == 0 || len(southHash) == 0 {
			t.Errorf("Geohash should not be empty for pole coordinates")
		}

		// Distance between poles should be approximately half Earth's circumference
		distance := northPole.DistanceFrom(southPole)
		expectedDistance := PI * 2 * EarthRadius / 2 // Approximate distance between poles in meters
		if distance < expectedDistance*0.99 || distance > expectedDistance*1.01 {
			t.Errorf("Distance between poles %f is not close to expected %f", distance, expectedDistance)
		}
	})

	t.Run("Antipodal points", func(t *testing.T) {
		point1 := NewCoordinateFromFloat64(40.7128, -74.0060)  // NYC
		point2 := NewCoordinateFromFloat64(-40.7128, 105.9940) // Antipodal to NYC

		distance := point1.DistanceFrom(point2)
		// Should be approximately half Earth's circumference
		expectedDistance := PI * 2 * EarthRadius / 2

		if distance < expectedDistance*0.98 || distance > expectedDistance*1.02 {
			t.Errorf("Distance between antipodal points %f is not close to expected %f", distance, expectedDistance)
		}
	})
}

func BenchmarkDistanceFrom(b *testing.B) {
	coord1 := NewCoordinateFromFloat64(40.7128, -74.0060)
	coord2 := NewCoordinateFromFloat64(34.0522, -118.2437)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		coord1.DistanceFrom(coord2)
	}
}

func BenchmarkCalGeoHash(b *testing.B) {
	coord := NewCoordinateFromFloat64(40.7128, -74.0060)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		coord.CalGeoHash()
	}
}
