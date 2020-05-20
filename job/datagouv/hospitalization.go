package datagouv

import (
	"io"
	"strconv"
	"time"

	"github.com/thetreep/covidtracker"
)

var _ covidtracker.HospService = &Service{}

func (s *Service) RefreshHospitalization() ([]*covidtracker.Hospitalization, error) {

	//TODO add limits to avoid duplicate
	//header
	const (
		dep = iota
		sexe
		jour
		hosp
		rea
		rad
		dc
	)

	url := "https://www.data.gouv.fr/fr/datasets/r/6fadff46-9efd-4c53-942a-54aca783c30c"
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

		entry := &covidtracker.Hospitalization{
			Sex: line[sexe],
		}

		entry.Department, err = atoi(line[dep])
		if s.handleParsingErr(err, "hospitaization", "dep") != nil {
			continue
		}
		entry.Count, err = atoi(line[hosp])
		if s.handleParsingErr(err, "hospitaization", "hosp") != nil {
			continue
		}
		entry.CriticalCount, err = atoi(line[rea])
		if s.handleParsingErr(err, "hospitaization", "rea") != nil {
			continue
		}
		entry.ReturnHomeCount, err = atoi(line[rad])
		if s.handleParsingErr(err, "hospitaization", "rad") != nil {
			continue
		}
		entry.DeathCount, err = atoi(line[dc])
		if s.handleParsingErr(err, "hospitaization", "dc") != nil {
			continue
		}
		entry.NoticeDate, err = time.Parse("2006-02-01", line[jour])
		if s.handleParsingErr(err, "hospitaization", "jour") != nil {
			continue
		}

		result = append(result, entry)
	}

	return result, nil
}