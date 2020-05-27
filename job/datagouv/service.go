package datagouv

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
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

	DatagouvBase = "https://www.data.gouv.fr"
)

func (s *Service) newReader(r io.Reader, res Resource) (*csv.Reader, error) {

	reader := csv.NewReader(r)
	switch res {
	case EmergencyURL, IndicatorURL:
		reader.Comma = ','
	case CaseURL, HospitalizationURL, ScreeningURL:
		reader.Comma = ';'
	default:
		return nil, fmt.Errorf("unsupported resource %q", res)
	}

	return reader, nil
}

func (s *Service) GetCSV(res Resource) (*csv.Reader, func() error, error) {
	// Get the data
	resp, err := http.Get(s.BasePath + string(res))
	if err != nil {
		return nil, nil, err
	}

	r, err := s.newReader(resp.Body, res)
	if err != nil {
		return nil, nil, err
	}

	return r, resp.Body.Close, nil
}

func (s *Service) handleParsingErr(err error, name, col string) error {
	if err != nil {
		//TODO use logger of service
		fmt.Printf("%s : cannot parse %q column (%v)\n", name, col, err)
	}
	return err
}
