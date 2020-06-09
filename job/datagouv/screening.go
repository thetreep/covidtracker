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

	reader, close, err := s.GetCSV(ScreeningID)
	if err != nil {
		return nil, err
	}
	defer close()

	var (
		result []*covidtracker.Screening
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

		entry := &covidtracker.Screening{
			AgeGroup: covidtracker.AgeGroup(line[clageCovid]),
		}

		entry.Department, err = atoi(line[dep])
		if s.handleParsingErr(err, "screening", "dep") != nil {
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
		entry.NoticeDate, err = time.Parse("2006-01-02", line[jour])
		if s.handleParsingErr(err, "screening", "jour") != nil {
			continue
		}

		result = append(result, entry)
	}

	return result, nil
}
