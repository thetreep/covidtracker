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

package job

import (
	"strings"
	"testing"
	"time"

	"github.com/thetreep/covidtracker"
	"github.com/thetreep/covidtracker/mock"
	"github.com/thetreep/toolbox/convert"
	"github.com/thetreep/toolbox/test"
)

func TestComputeRisk(t *testing.T) {
	paris := &covidtracker.Geo{Properties: covidtracker.Properties{DepCode: convert.StrP("75")}}
	lyon := &covidtracker.Geo{Properties: covidtracker.Properties{DepCode: convert.StrP("69")}}
	dep, _ := time.Parse("02/01/2006 15:04:05", "18/03/2020 13:45:00")
	arrival := dep.Add(2*time.Hour + 10*time.Minute)

	parametersMock := mock.RiskParameters{
		GetDefaultFn: func() (*covidtracker.RiskParameters, error) {
			return &covidtracker.RiskParameters{
				IsDefault:                true,
				SewnMaskProtect:          0.5,
				SurgicalMaskProtect:      0.90,
				FFPXMaskProtect:          0.99,
				HydroAlcoholicGelProtect: 0.99,
				Parameters: []*covidtracker.RiskParameter{
					{
						Scope:                 covidtracker.ParameterScope{Transportation: covidtracker.TGV, Duration: covidtracker.Normal},
						NbContact:             4,
						MaskProtectContact:    0.2,
						GelProtectContact:     0.9,
						ProbaContagionContact: 0.9,
						NbDirect:              10,
						ProbaContagionDirect:  0.7,
						MaskProtectDirect:     0.8,
						Advices:               []string{"restez chez vous !"},
						Minuses:               []string{"ne prenez pas le train"},
					},
				},
			}, nil
		},
	}

	emergencyMock := mock.EmergencyDAL{
		GetRangeFn: func(dep string, start, end time.Time) ([]*covidtracker.Emergency, error) {
			var res []*covidtracker.Emergency
			for i := 0; i < 14; i++ {
				emerg := &covidtracker.Emergency{Cov19SuspCount: 1000 * (i + 1)}
				res = append(res, emerg)
			}
			return res, nil
		},
	}

	tcases := []struct {
		name              string
		segments          []covidtracker.Segment
		protects          []covidtracker.Protection
		expectRisk        *covidtracker.Risk
		expectErrContains string
	}{
		{
			name:              "no segments",
			expectErrContains: "no segment found",
		},
		{
			name: "OK",
			segments: []covidtracker.Segment{
				{Origin: paris, Destination: lyon, Departure: dep, Arrival: arrival, Transportation: covidtracker.TGV},
			},
			protects: []covidtracker.Protection{{Type: covidtracker.Gel}, {Type: covidtracker.MaskFFPX}, {Type: covidtracker.MaskSewn}},
			expectRisk: &covidtracker.Risk{
				NoticeDate:      dep,
				ConfidenceLevel: 0.9580939763332113,
				RiskLevel:       0.04190602366678872,
				DisplayedRisk:   0.35953011833394366,
				BySegments: []covidtracker.RiskSegment{{
					Segment: &covidtracker.Segment{
						Origin:         paris,
						Destination:    lyon,
						Departure:      dep,
						Arrival:        arrival,
						Transportation: covidtracker.TGV,
					},
					ConfidenceLevel: 0.9580939763332113,
					RiskLevel:       0.04190602366678872,
					Report: covidtracker.Report{
						Advices: []covidtracker.Statement{{Value: "restez chez vous !", Category: "tgv"}},
						Minuses: []covidtracker.Statement{{Value: "ne prenez pas le train", Category: "tgv"}},
					},
				}},
				Report: covidtracker.Report{
					Advices: []covidtracker.Statement{{Value: "restez chez vous !", Category: "tgv"}},
					Pluses:  []covidtracker.Statement{{Value: "Vous portez un masque", Category: "mask"}, {Value: "Vous utilisez du gel hydroalcoolique", Category: "gel"}},
					Minuses: []covidtracker.Statement{{Value: "ne prenez pas le train", Category: "tgv"}},
				},
			},
		},
	}

	j := &RiskJob{job: &Job{RiskParametersDAL: &parametersMock, EmergencyDAL: &emergencyMock, Now: func() time.Time { return dep }}}
	for _, tcase := range tcases {
		r, err := j.ComputeRisk(tcase.segments, tcase.protects)
		if tcase.expectErrContains != "" {
			if err == nil {
				t.Fatalf("%s: expect error to contain %q but got nil", tcase.name, tcase.expectErrContains)
			}
			if got, want := err.Error(), tcase.expectErrContains; !strings.Contains(got, want) {
				t.Fatalf("%s: got %q, want %q", tcase.name, got, want)
			}
		} else {
			if err != nil {
				t.Fatalf("%s: %s", tcase.name, err)
			}
			test.Compare(t, r, tcase.expectRisk, tcase.name)
		}
	}
}

func TestYOnLine(t *testing.T) {
	tcases := []struct {
		x       float64
		x1, y1  float64
		x2, y2  float64
		expectY float64
	}{
		{x: 0, x1: 0, y1: 0, x2: 4, y2: 4, expectY: 0},
		{x: 1, x1: 0, y1: 0, x2: 4, y2: 4, expectY: 1},
		{x: 0, x1: 0, y1: 5, x2: 4, y2: 4, expectY: 5},
		{x: 4, x1: 0, y1: 5, x2: 4, y2: 4, expectY: 4},
		{x: 4, x1: 2, y1: 6, x2: 12, y2: 11, expectY: 7},
	}
	for i, tcase := range tcases {
		y := yOnLine(tcase.x, tcase.x1, tcase.y1, tcase.x2, tcase.y2)
		if got, want := y, tcase.expectY; got != want {
			t.Fatalf("%d: got %f, want %f", i+1, got, want)
		}
	}
}
