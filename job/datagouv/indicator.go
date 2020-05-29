package datagouv

import (
	"io"
	"strconv"
	"time"

	"github.com/thetreep/covidtracker"
)

var _ covidtracker.IndicService = &Service{}

func (s *Service) RefreshIndicator() ([]*covidtracker.Indicator, error) {

	//Header of csv emergency file
	const (
		extractDate = iota
		departement
		indicSynthese
	)

	reader, close, err := s.GetCSV(IndicatorURL)
	if err != nil {
		return nil, err
	}
	defer close()

	var (
		result []*covidtracker.Indicator
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

		entry := &covidtracker.Indicator{
			Color: line[indicSynthese],
		}
		entry.Department, err = atoi(line[departement])
		if s.handleParsingErr(err, "indicator", "dep") != nil {
			continue
		}
		entry.ExtractDate, err = time.Parse("2006-01-02", line[extractDate])
		if s.handleParsingErr(err, "indicator", "extractDate") != nil {
			continue
		}

		result = append(result, entry)
	}

	return result, nil
}
