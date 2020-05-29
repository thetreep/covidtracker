package covidtracker

import (
	"fmt"
	"regexp"
)

type Geo struct {
	Properties Properties `json:"properties"`
	Type       string     `json:"type"`
	Geometry   Geometry   `json:"geometry"`
}

type Properties struct {
	GeoCoding GeoCoding `json:"geocoding"`
}

type Geometry struct {
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}

type GeoCoding struct {
	Type        string      `json:"type"`
	Accuracy    int         `json:"accuracy"`
	Label       string      `json:"label"`
	Name        string      `json:"name"`
	HouseNumber string      `json:"housenumber"`
	Street      string      `json:"street"`
	Locality    string      `json:"locality"`
	PostCode    string      `json:"postcode"`
	City        string      `json:"city"`
	District    *string     `json:"district,omitempty"`
	County      *string     `json:"county,omitempty"`
	State       *string     `json:"state,omitempty"`
	Country     string      `json:"country,omitempty"`
	Admin       AdminLevels `json:"admin"`
	Geohash     string      `json:"geohash"`
}

type AdminLevels struct {
	Level2 string `json:"level2"`
	Level4 string `json:"level4"`
	Level6 string `json:"level6"`
}

func (g Geo) Check() error {

	regPostCode := `^(([0-8][0-9])|(9[0-5])|(2[ab]))[0-9]{3}$`

	match, err := regexp.Match(regPostCode, []byte(g.Properties.GeoCoding.PostCode))
	if err != nil {
		return err
	}

	if !match {
		return fmt.Errorf("postal code %q is invalid", g.Properties.GeoCoding.PostCode)
	}

	return nil
}
