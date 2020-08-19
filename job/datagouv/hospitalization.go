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

var _ covidtracker.HospService = &Service{}

func (s *Service) RefreshHospitalization() ([]*covidtracker.Hospitalization, error) {
	s.log.Debug(s.Ctx, "refreshing hospitalization data...")

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

		entry := &covidtracker.Hospitalization{
			Department: line[dep],
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
		if err != nil {
			//there are two format possible in this file...
			entry.Date, err = time.Parse("02/01/2006", line[jour])
			if s.handleParsingErr(err, "hospitalization", "jour") != nil {
				continue
			}
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

	s.log.Debug(s.Ctx, "got %d hospitalization entries !", len(result))

	return result, nil
}
