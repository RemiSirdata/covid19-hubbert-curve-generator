package pkg

import (
	"fmt"
	"strings"
)

const (
	dateFormat = "2006-01-02"
)

type RawData struct {
}

func (r *RawData) GetCountry(country string) (RawCountryData, error) {
	switch strings.ToLower(country) {
	case "france":
		return newRawFranceData()
	}
	return nil, fmt.Errorf("country not available")
}
