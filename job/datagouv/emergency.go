/*
	This file is part of covidtracker also known as EviteCovid .

    Copyright 2020 the Treep

    covdtracker is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    covidtracker is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with covidtracker.  If not, see <https://www.gnu.org/licenses/>.
*/

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
	s.log.Debug(s.Ctx, "refreshing emergency data...")

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

		entry := &covidtracker.Emergency{
			Department: line[dep],
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

	s.log.Debug(s.Ctx, "got %d emergency entries !", len(result))

	return result, nil
}
