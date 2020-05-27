package datagouv

import (
	"io"
	"strconv"
	"time"

	"github.com/thetreep/covidtracker"
)

var _ covidtracker.CaseService = &Service{}

func (s *Service) RefreshCase() ([]*covidtracker.Case, error) {
	//header
	const (
		dep = iota
		jour
		nb
	)

	reader, err := s.getCSV(CaseURL)
	if err != nil {
		return nil, err
	}

	var (
		result []*covidtracker.Case
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

		entry := &covidtracker.Case{}

		entry.Department, err = atoi(line[dep])
		if s.handleParsingErr(err, "covid_case", "nbTest") != nil {
			continue
		}
		entry.HospServiceCountRelated, err = atoi(line[nb])
		if s.handleParsingErr(err, "covid_case", "txPos") != nil {
			continue
		}
		entry.NoticeDate, err = time.Parse("2006-02-01", line[jour])
		if s.handleParsingErr(err, "covid_case", "jour") != nil {
			continue
		}

		result = append(result, entry)
	}
	return result, nil
}
