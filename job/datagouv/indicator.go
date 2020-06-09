package datagouv

import (
	"io"
	"sort"
	"time"

	"github.com/thetreep/covidtracker"
)

var _ covidtracker.IndicService = &Service{}

func (s *Service) RefreshIndicator() ([]*covidtracker.Indicator, error) {
	s.log.Debug(s.Ctx, "refreshing indicator data...")

	//Header of csv indicator file
	const (
		extractDate = iota
		departement
		depName
		region
		indicSynthese
	)

	reader, close, err := s.GetCSV(IndicatorID)
	if err != nil {
		return nil, err
	}
	defer close()

	var (
		result      []*covidtracker.Indicator
		resultByKey = make(map[string]*covidtracker.Indicator)
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
			Color:      line[indicSynthese],
			Department: line[departement],
		}

		entry.ExtractDate, err = time.Parse("2006-01-02", line[extractDate])
		if s.handleParsingErr(err, "indicator", "extractDate") != nil {
			continue
		}

		//avoid duplicate
		k := line[extractDate] + "_" + line[departement]
		resultByKey[k] = entry

	}

	for _, e := range resultByKey {
		result = append(result, e)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].ExtractDate.Equal(result[j].ExtractDate) {
			return result[i].Department < result[j].Department
		}
		return result[i].ExtractDate.After(result[j].ExtractDate)
	})

	s.log.Debug(s.Ctx, "got %d screening data !", len(result))

	return result, nil
}
