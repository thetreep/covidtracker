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
				DisplayedRisk:   0.42057033991661624,
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
