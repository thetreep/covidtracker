package covidtracker

import "fmt"

type HotelID string

type Hotel struct {
	ID            HotelID  `bson:"_id" json:"id"`
	Name          string   `bson:"name" json:"name"`
	Address       string   `bson:"address" json:"address"`
	City          string   `bson:"city" json:"city"`
	ZipCode       string   `bson:"zipCode" json:"zipCode"`
	Country       string   `bson:"country" json:"country"`
	ImageURL      string   `bson:"imageUrl" json:"imageUrl"`
	SanitaryInfos []string `bson:"sanitaryInfos" json:"sanitaryInfos"`
	SanitaryNote  float64  `bson:"sanitaryNote" json:"sanitaryNote"`
	SanitaryNorm  string   `bson:"sanitaryNorm" json:"sanitaryNorm"`
}

func (h *Hotel) Dep() (string, error) {
	if len(h.ZipCode) > 2 {
		return h.ZipCode[:2], nil
	}

	return "", fmt.Errorf("department missing")
}

type HotelDAL interface {
	Get(id HotelID) (*Hotel, error)
	Insert(hotels []*Hotel) ([]*Hotel, error)
}

type HotelJob interface {
	HotelsByPrefix(city string, prefix string) ([]*Hotel, error)
}
