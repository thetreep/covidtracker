package datagouv

import (
	"io"
	"sort"
	"strconv"
	"time"

	"github.com/thetreep/covidtracker"
)

var _ covidtracker.CaseService = &Service{}

func (s *Service) RefreshCase() ([]*covidtracker.Case, error) {
	s.log.Debug(s.Ctx, "refreshing case data...")

	//header
	const (
		dep = iota
		jour
		nb
	)

	reader, close, err := s.GetCSV(CaseID)
	if err != nil {
		return nil, err
	}
	defer close()

	var (
		result      []*covidtracker.Case
		resultByKey = make(map[string]*covidtracker.Case)
		atoi        = strconv.Atoi
	)

	reader.Read() //ignore first line (columns names)
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		entry := &covidtracker.Case{
			Department: line[dep],
		}

		entry.HospServiceCountRelated, err = atoi(line[nb])
		if s.handleParsingErr(err, "covid_case", "txPos") != nil {
			continue
		}
		entry.NoticeDate, err = time.Parse("2006-01-02", line[jour])
		if s.handleParsingErr(err, "covid_case", "jour") != nil {
			continue
		}

		//avoid duplicate
		k := line[jour] + "_" + line[dep]
		resultByKey[k] = entry

	}

	for _, e := range resultByKey {
		result = append(result, e)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].NoticeDate.Equal(result[j].NoticeDate) {
			return result[i].Department < result[j].Department
		}
		return result[i].NoticeDate.After(result[j].NoticeDate)
	})

	s.log.Debug(s.Ctx, "got %d case entries !", len(result))

	return result, nil
}
