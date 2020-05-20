package datagouv

import (
	"io"
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
	url := "https://www.data.gouv.fr/fr/datasets/r/eceb9fb4-3ebc-4da3-828d-f5939712600a"
	reader, err := s.getCSV(url)
	if err != nil {
		return nil, err
	}

	var (
		result []*covidtracker.Emergency
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

		entry := &covidtracker.Emergency{
			AgeGroup: line[sursaudClAgeCorona],
		}
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
		entry.NoticeDate, err = time.Parse("2006-02-01", line[dateDePassage])
		if s.handleParsingErr(err, "emergency", "dateDePassage") != nil {
			continue
		}

		result = append(result, entry)
	}

	return result, nil
}
