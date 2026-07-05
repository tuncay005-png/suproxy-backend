package server

// Location value object
type Location struct {
	Country     string
	City        string
	CountryCode string
	Latitude    float64
	Longitude   float64
}

func NewLocation(country, city, countryCode string, lat, lon float64) Location {
	return Location{
		Country:     country,
		City:        city,
		CountryCode: countryCode,
		Latitude:    lat,
		Longitude:   lon,
	}
}

// Capacity value object
type Capacity struct {
	MaxUsers       int
	CurrentUsers   int
	MaxBandwidth   int64
	UsedBandwidth  int64
}

func NewCapacity(maxUsers int, maxBandwidth int64) Capacity {
	return Capacity{
		MaxUsers:      maxUsers,
		CurrentUsers:  0,
		MaxBandwidth:  maxBandwidth,
		UsedBandwidth: 0,
	}
}

func (c Capacity) IsFull() bool {
	return c.CurrentUsers >= c.MaxUsers
}

func (c Capacity) AvailableSlots() int {
	return c.MaxUsers - c.CurrentUsers
}

func (c Capacity) UsagePercentage() float64 {
	if c.MaxUsers == 0 {
		return 0
	}
	return float64(c.CurrentUsers) / float64(c.MaxUsers) * 100
}

func (c Capacity) BandwidthUsagePercentage() float64 {
	if c.MaxBandwidth == 0 {
		return 0
	}
	return float64(c.UsedBandwidth) / float64(c.MaxBandwidth) * 100
}

// Status enum
type Status string

const (
	StatusActive      Status = "active"
	StatusInactive    Status = "inactive"
	StatusMaintenance Status = "maintenance"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusInactive, StatusMaintenance:
		return true
	}
	return false
}
