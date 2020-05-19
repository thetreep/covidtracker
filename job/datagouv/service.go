package datagouv

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
)

type Service struct {
	Ctx context.Context

	//TODO define logger
}

func (s *Service) getCSV(url string) (*csv.Reader, error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return csv.NewReader(resp.Body), nil
}

func (s *Service) handleParsingErr(err error, name, col string) error {
	if err != nil {
		//TODO use logger of service
		fmt.Println("%s : cannot parse %q column (%v)", name, col, err)
	}
	return err
}
