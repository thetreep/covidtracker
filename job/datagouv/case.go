package datagouv

import (
	"io"
	"strconv"
	"time"

	"github.com/thetreep/covidtracker"
)

var _ covidtracker.CaseService = &Service{}

func (s *Service) RefreshCase() ([]*covidtracker.Case, error) {
	//TODO add limits to avoid duplicate

	//header
	const (
		dep = iota
		jour
		nb
	)

	url := "https://www.data.gouv.fr/fr/datasets/r/b4ea7b4b-b7d1-4885-a099-71852291ff20"
	reader, err := s.getCSV(url)
	if err != nil {
		return nil, err
	}

	var (
		result []*covidtracker.Hospitalization
		atoi   = strconv.Atoi
	)

	reader.Read() //ignore first line (columns names)
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		entry := &covidtracker.Hospitalization{}

		entry.Department, err = atoi(line[dep])
		if s.handleParsingErr(err, "covid_case", "nbTest") != nil {
			continue
		}
		entry.HospServiceCountRelated, err = atoi(line[nb])
		if s.handleParsingErr(err, "screening", "txPos") != nil {
			continue
		}
		entry.NoticeDate, err = time.Parse("2006-02-01", line[jour])
		if s.handleParsingErr(err, "screening", "jour") != nil {
			continue
		}

		result = append(result, entry)
	}
}
