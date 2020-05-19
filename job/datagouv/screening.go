package datagouv

import (
	"io"
	"strconv"
	"time"

	"github.com/thetreep/covidtracker"
)

var _ covidtracker.ScreeningService = &Service{}

func (s *Service) RefreshScreening() ([]*covidtracker.Screening, error) {
	//TODO add limits to avoid duplicate

	//header
	const (
		dep = iota
		jour
		clageCovid
		nbTest
		nbPos
		txPos
		nbTestH
		nbPosH
		nbTestF
		nbPosF
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

		entry := &covidtracker.Hospitalization{
			AgeGroup: line[clageCovid],
		}

		entry.Department, err = atoi(line[dep])
		if s.handleParsingErr(err, "emergency", "dep") != nil {
			continue
		}
		entry.Count, err = atoi(line[nbTest])
		if s.handleParsingErr(err, "screening", "nbTest") != nil {
			continue
		}
		entry.PositiveCount, err = atoi(line[nbPos])
		if s.handleParsingErr(err, "screening", "nbPos") != nil {
			continue
		}
		entry.PositiveRate, err = atoi(line[txPos])
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
