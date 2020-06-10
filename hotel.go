package covidtracker

type Hotel struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Address       string   `json:"address"`
	City          string   `json:"city"`
	ZipCode       string   `json:"zip_code"`
	Country       string   `json:"country"`
	ImageURL      string   `json:"ImageUrl"`
	SanitaryInfos []string `json:"sanitary_infos"`
	SanitaryNote  float64  `json:"sanitary_note"`
	SanitaryNorm  string   `json:"sanitary_norm"`
}

type HotelJob interface {
	HotelsByPrefix(prefix string) ([]*Hotel, error)
}
