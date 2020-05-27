package datagouv

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"path"
)

type Service struct {
	Ctx context.Context

	BasePath string
	//TODO define logger
}

type Resource string

const (
	EmergencyURL       Resource = "/fr/datasets/r/eceb9fb4-3ebc-4da3-828d-f5939712600a"
	CaseURL            Resource = "/fr/datasets/r/41b9bd2a-b5b6-4271-8878-e45a8902ef00"
	HospitalizationURL Resource = "/fr/datasets/r/6fadff46-9efd-4c53-942a-54aca783c30c"
	ScreeningURL       Resource = "/fr/datasets/r/b4ea7b4b-b7d1-4885-a099-71852291ff20"
	IndicatorURL       Resource = "/fr/datasets/r/01151af0-3209-4e89-94ab-9b319001c159"

	Datagouv = "https://www.data.gouv.fr"
)

func (s *Service) getCSV(url Resource) (*csv.Reader, error) {
	// Get the data
	resp, err := http.Get(path.Join(s.BasePath, url))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return csv.NewReader(resp.Body), nil
}

func (s *Service) handleParsingErr(err error, name, col string) error {
	if err != nil {
		//TODO use logger of service
		fmt.Printf("%s : cannot parse %q column (%v)", name, col, err)
	}
	return err
}
