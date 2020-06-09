package datagouv

import (
	"io"
	"sort"
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

	reader, close, err := s.GetCSV(HospitalizationID)
	if err != nil {
		return nil, err
	}
	defer close()

	var (
		result      []*covidtracker.Hospitalization
		resultByKey = make(map[string]*covidtracker.Hospitalization)
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

		entry := &covidtracker.Hospitalization{}

		entry.Department, err = atoi(line[dep])
		if s.handleParsingErr(err, "hospitalization", "dep") != nil {
			continue
		}
		entry.Count, err = atoi(line[hosp])
		if s.handleParsingErr(err, "hospitalization", "hosp") != nil {
			continue
		}
		entry.CriticalCount, err = atoi(line[rea])
		if s.handleParsingErr(err, "hospitalization", "rea") != nil {
			continue
		}
		entry.ReturnHomeCount, err = atoi(line[rad])
		if s.handleParsingErr(err, "hospitalization", "rad") != nil {
			continue
		}
		entry.DeathCount, err = atoi(line[dc])
		if s.handleParsingErr(err, "hospitalization", "dc") != nil {
			continue
		}
		entry.Date, err = time.Parse("2006-01-02", line[jour])
		if s.handleParsingErr(err, "hospitalization", "jour") != nil {
			continue
		}

		k := line[jour] + "_" + line[dep]
		if _, ok := resultByKey[k]; ok {
			resultByKey[k].Count += entry.Count
			resultByKey[k].DeathCount += entry.DeathCount
			resultByKey[k].CriticalCount += entry.CriticalCount
			resultByKey[k].ReturnHomeCount += entry.ReturnHomeCount
		} else {
			resultByKey[k] = entry
		}
	}

	for _, e := range resultByKey {
		result = append(result, e)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Date.Equal(result[j].Date) {
			return result[i].Department < result[j].Department
		}
		return result[i].Date.After(result[j].Date)
	})

	return result, nil
}
