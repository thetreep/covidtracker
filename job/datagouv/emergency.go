package datagouv

import (
	"io"
	"sort"
	"strconv"
	"time"

	"github.com/thetreep/covidtracker"
)

var _ covidtracker.EmergencyService = &Service{}

func (s *Service) RefreshEmergency() ([]*covidtracker.Emergency, error) {

	//Header of csv emergency file
	const (
		dep = iota
		dateDePassage
		sursaudClAgeCorona
		nbrePassCorona
		nbrePassTot
		nbreHospitCorona
		nbrePassCoronaH
		nbrePassCoronaF
		nbrePassTotH
		nbrePassTotF
		nbreHospitCoronaH
		nbreHospitCoronaF
		nbreActeCorona
		nbreActeTot
		nbreActeCoronaH
		nbreActeCoronaF
		nbreActeTotH
		nbreActeTotF
	)

	//TODO add limits to avoid duplicate
	reader, close, err := s.GetCSV(EmergencyID)
	if err != nil {
		return nil, err
	}
	defer close()

	var (
		result      []*covidtracker.Emergency
		resultByKey = make(map[string]*covidtracker.Emergency)
		atoi        = func(s string) (int, error) {
			if s == "" {
				return 0, nil
			}
			return strconv.Atoi(s)
		}
	)

	reader.Read() //ignore first line (columns names)
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		entry := &covidtracker.Emergency{}

		entry.Department, err = atoi(line[dep])
		if s.handleParsingErr(err, "emergency", "dep") != nil {
			continue
		}
		entry.Count, err = atoi(line[nbrePassTot])
		if s.handleParsingErr(err, "emergency", "nbrePassTot") != nil {
			continue
		}
		entry.Cov19SuspCount, err = atoi(line[nbrePassCorona])
		if s.handleParsingErr(err, "emergency", "nbrePassCorona") != nil {
			continue
		}
		entry.Cov19SuspHosp, err = atoi(line[nbreHospitCorona])
		if s.handleParsingErr(err, "emergency", "nbreHospitCorona") != nil {
			continue
		}
		entry.TotalSOSMedAct, err = atoi(line[nbreActeTot])
		if s.handleParsingErr(err, "emergency", "nbreActeTot") != nil {
			continue
		}
		entry.SOSMedCov19SuspAct, err = atoi(line[nbreActeCorona])
		if s.handleParsingErr(err, "emergency", "nbreActeCorona") != nil {
			continue
		}
		entry.PassageDate, err = time.Parse("2006-01-02", line[dateDePassage])
		if s.handleParsingErr(err, "emergency", "dateDePassage") != nil {
			continue
		}

		k := line[dateDePassage] + "_" + line[dep]
		if _, ok := resultByKey[k]; ok {
			resultByKey[k].SOSMedCov19SuspAct += entry.SOSMedCov19SuspAct
			resultByKey[k].TotalSOSMedAct += entry.TotalSOSMedAct
			resultByKey[k].Cov19SuspHosp += entry.Cov19SuspHosp
			resultByKey[k].Cov19SuspCount += entry.Cov19SuspCount
			resultByKey[k].Count += entry.Count
		} else {
			resultByKey[k] = entry
		}
	}

	for _, e := range resultByKey {
		result = append(result, e)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].PassageDate.Equal(result[j].PassageDate) {
			return result[i].Department < result[j].Department
		}
		return result[i].PassageDate.After(result[j].PassageDate)
	})

	return result, nil
}
