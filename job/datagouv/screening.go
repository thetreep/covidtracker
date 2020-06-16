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

var _ covidtracker.ScreeningService = &Service{}

func (s *Service) RefreshScreening() ([]*covidtracker.Screening, error) {
	s.log.Debug(s.Ctx, "refreshing screening data...")

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
		result      []*covidtracker.Screening
		resultByKey = make(map[string]*covidtracker.Screening)
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

		entry := &covidtracker.Screening{
			Department: line[dep],
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

		k := line[jour] + "_" + line[dep]
		if _, ok := resultByKey[k]; ok {
			resultByKey[k].PositiveRate += entry.PositiveRate
			resultByKey[k].PositiveCount += entry.PositiveCount
			resultByKey[k].Count += entry.Count
		} else {
			resultByKey[k] = entry
		}
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

	s.log.Debug(s.Ctx, "got %d screening entries !", len(result))

	return result, nil
}
