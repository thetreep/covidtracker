package covidtracker

type HotelID string

type Hotel struct {
	ID            HotelID  `bson:"_id" json:"id"`
	Name          string   `bson:"name" json:"name"`
	Address       string   `bson:"address" json:"address"`
	City          string   `bson:"city" json:"city"`
	ZipCode       string   `bson:"zip_code" json:"zip_code"`
	Country       string   `bson:"country" json:"country"`
	ImageURL      string   `bson:"image_url" json:"image_url"`
	SanitaryInfos []string `bson:"sanitary_infos" json:"sanitary_infos"`
	SanitaryNote  float64  `bson:"sanitary_note" json:"sanitary_note"`
	SanitaryNorm  string   `bson:"sanitary_norm" json:"sanitary_norm"`
}

type HotelDAL interface {
	Get(id HotelID) (*Hotel, error)
	Insert(hotels []*Hotel) ([]*Hotel, error)
}

type HotelJob interface {
	HotelsByPrefix(prefix string) ([]*Hotel, error)
}
