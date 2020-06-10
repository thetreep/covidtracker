package graphql_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/thetreep/covidtracker"
	"github.com/thetreep/covidtracker/http/graphql"
	"github.com/thetreep/covidtracker/mock"

	"github.com/thetreep/toolbox/convert"
	"github.com/thetreep/toolbox/test"
)

func TestEstimate(t *testing.T) {
	ctx := context.Background()

	risk := mock.Risk{}
	h, err := graphql.NewHandler(&graphql.RiskHandler{DAL: &risk, Job: &risk})
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(h)
	defer server.Close()

	client := NewClient(ctx, server.URL)

	tplFull := `
	query ($segs: [segmentIn], $prots: [protectionIn], $hots: [hotelIn]) {
		risk(
			segments:$segs,
			protections:$prots,
			hotels:$hots
		) {
			bySegments {
				riskLevel
				segment {
					origin { properties { geocoding {postcode,city,locality}}}
					destination { properties { geocoding {postcode,city,locality}}}
				}
			}
			riskLevel
			confidenceLevel
			report {
				minuses {value}
				pluses {value}
				advices {value}
			}
		}
	}`

	tplReport := `
	query ($segs: [segmentIn], $prots: [protectionIn], $hots: [hotelIn]) {
		risk(
			segments:$segs,
			protections:$prots,
			hotels:$hots
		) {
			report {
				minuses {value}
				pluses {value}
				advices {value}
			}
		}
	}`

	tplRisk := `
	query ($segs: [segmentIn], $prots: [protectionIn], $hots: [hotelIn]) {
		risk(
			segments:$segs,
			protections:$prots,
			hotels:$hots
		) {
			riskLevel
		}
	}`

	t.Run("no segment provided", func(t *testing.T) {
		got, err := client.Do(tplFull, map[string]interface{}{
			"segs":  []covidtracker.Segment{},
			"prots": []covidtracker.Protection{},
			"hots":  []graphql.HotelInput{},
		})
		expected := &gqlResp{
			Data:   map[string]interface{}{"risk": nil},
			Errors: []gqlErr{{Message: "at least one `segment` or `hotel` is mandatory"}},
		}
		test.Compare(t, err.Error(), "at least one `segment` or `hotel` is mandatory", "unexpected error")
		test.Compare(t, got, expected, "unexpected result")
		test.Compare(t, risk.ComputeRiskInvoked, false, "estimate invokation unexpected")
		test.Compare(t, risk.InsertInvoked, false, "insert invokation unexpected")
	})

	t.Run("computation error", func(t *testing.T) {
		risk.ComputeRiskFn = func(segs []covidtracker.Segment, protects []covidtracker.Protection) (*covidtracker.Risk, error) {
			return nil, fmt.Errorf("computation error")
		}

		got, err := client.Do(tplFull, map[string]interface{}{
			"segs":  []covidtracker.Segment{{Origin: &mock.Paris, Destination: &mock.Bordeaux}},
			"prots": []covidtracker.Protection{},
			"hots":  []graphql.HotelInput{},
		})
		expected := &gqlResp{
			Data:   map[string]interface{}{"risk": nil},
			Errors: []gqlErr{{Message: "computation error"}},
		}
		test.Compare(t, err.Error(), "computation error", "unexpected error")
		test.Compare(t, got, expected, "unexpected result")
		test.Compare(t, risk.ComputeRiskInvoked, true, "estimate invokation is expected")
		test.Compare(t, risk.InsertInvoked, false, "insert invokation unexpected")
	})

	t.Run("computation OK", func(t *testing.T) {

		risk.ComputeRiskFn = func(segs []covidtracker.Segment, protects []covidtracker.Protection) (*covidtracker.Risk, error) {
			r := &covidtracker.Risk{
				RiskLevel:       0.5,
				ConfidenceLevel: 0.5,
				Report: covidtracker.Report{
					Minuses: []covidtracker.Statement{{Value: "it's not good"}},
					Pluses:  []covidtracker.Statement{{Value: "it's pretty good"}},
					Advices: []covidtracker.Statement{{Value: "a wise advice"}},
				},
			}
			for i := range segs {
				rs := covidtracker.RiskSegment{
					ID:              covidtracker.RiskSegID(fmt.Sprint(i + 1)),
					Segment:         &segs[i],
					RiskLevel:       .5,
					ConfidenceLevel: .5,
				}
				if segs[i].HotelID != nil {
					r.RiskLevel = 0.6
					r.ConfidenceLevel = 0.4
					rs.RiskLevel = .7
					rs.ConfidenceLevel = .3
				}
				r.BySegments = append(r.BySegments, rs)
			}
			return r, nil
		}

		db := []*covidtracker.Risk{}
		risk.InsertFn = func(risks ...*covidtracker.Risk) error {
			for i, r := range risks {
				r.ID = covidtracker.RiskID(fmt.Sprint(i + 1))
				db = append(db, r)
			}
			return nil
		}

		tcases := map[string]struct {
			segs   []covidtracker.Segment
			hots   []graphql.HotelInput
			tpl    string
			expRaw []byte
			expDB  []*covidtracker.Risk
		}{
			"full query": {
				segs: []covidtracker.Segment{{Origin: &mock.Paris, Destination: &mock.Bordeaux}},
				hots: []graphql.HotelInput{},
				tpl:  tplFull,
				expRaw: []byte(`{
					"data": {
						"risk": {
							"riskLevel":       0.5,
							"confidenceLevel": 0.5,
							"bySegments": [{
								"segment": {
									"origin": {
										"properties": {
											"geocoding": {
												"postcode":"75015",
												"city":"Paris 15",
												"locality":"Paris"
											}
										}
									},
									"destination": {
										"properties": {
											"geocoding": {
												"postcode":"33000",
												"city":"Bordeaux",
												"locality":"Bordeaux"
											}
										}
									}
								},
								"riskLevel":       0.5
							}],
							"report": {
								"minuses": [{"value":"it's not good"}],
								"pluses": [{"value":"it's pretty good"}],
								"advices": [{"value":"a wise advice"}]
							}
						}
					}
				}`),
				expDB: []*covidtracker.Risk{{
					ID:              "1",
					RiskLevel:       .5,
					ConfidenceLevel: .5,
					BySegments: []covidtracker.RiskSegment{
						{
							ID:              covidtracker.RiskSegID("1"),
							Segment:         &covidtracker.Segment{Origin: &mock.Paris, Destination: &mock.Bordeaux},
							ConfidenceLevel: .5,
							RiskLevel:       .5,
						},
					},
					Report: covidtracker.Report{
						Advices: []covidtracker.Statement{{Value: "a wise advice"}},
						Pluses:  []covidtracker.Statement{{Value: "it's pretty good"}},
						Minuses: []covidtracker.Statement{{Value: "it's not good"}},
					},
				}},
			},
			"only risk": {
				segs: []covidtracker.Segment{{Origin: &mock.Paris, Destination: &mock.Bordeaux}},
				tpl:  tplRisk,
				expRaw: []byte(`{
					"data": {
						"risk": {
							"riskLevel":       0.5
						}
					}
				}`),
				expDB: []*covidtracker.Risk{{
					ID:              "1",
					RiskLevel:       .5,
					ConfidenceLevel: .5,
					BySegments: []covidtracker.RiskSegment{
						{
							ID:              covidtracker.RiskSegID("1"),
							Segment:         &covidtracker.Segment{Origin: &mock.Paris, Destination: &mock.Bordeaux},
							ConfidenceLevel: .5,
							RiskLevel:       .5,
						},
					},
					Report: covidtracker.Report{
						Advices: []covidtracker.Statement{{Value: "a wise advice"}},
						Pluses:  []covidtracker.Statement{{Value: "it's pretty good"}},
						Minuses: []covidtracker.Statement{{Value: "it's not good"}},
					},
				}},
			},
			"only risk with hotel": {
				segs: []covidtracker.Segment{{Origin: &mock.Paris, Destination: &mock.Bordeaux}},
				hots: []graphql.HotelInput{{ID: "an.hotel.id"}},
				tpl:  tplRisk,
				expRaw: []byte(`{
					"data": {
						"risk": {
							"riskLevel":       0.6
						}
					}
				}`),
				expDB: []*covidtracker.Risk{{
					ID:              "1",
					RiskLevel:       .6,
					ConfidenceLevel: .4,
					BySegments: []covidtracker.RiskSegment{
						{
							ID:              covidtracker.RiskSegID("1"),
							Segment:         &covidtracker.Segment{Origin: &mock.Paris, Destination: &mock.Bordeaux},
							ConfidenceLevel: .5,
							RiskLevel:       .5,
						},
						{
							ID:              covidtracker.RiskSegID("2"),
							Segment:         &covidtracker.Segment{HotelID: convert.StrP("an.hotel.id")},
							ConfidenceLevel: .3,
							RiskLevel:       .7,
						},
					},
					Report: covidtracker.Report{
						Advices: []covidtracker.Statement{{Value: "a wise advice"}},
						Pluses:  []covidtracker.Statement{{Value: "it's pretty good"}},
						Minuses: []covidtracker.Statement{{Value: "it's not good"}},
					},
				}}},
			"only report": {
				segs: []covidtracker.Segment{{Origin: &mock.Paris, Destination: &mock.Bordeaux}},
				hots: []graphql.HotelInput{},
				tpl:  tplReport,
				expRaw: []byte(`{
					"data": {
						"risk": {
							"report": {
								"minuses": [{"value":"it's not good"}],
								"pluses": [{"value":"it's pretty good"}],
								"advices": [{"value":"a wise advice"}]
							}
						}
					}
				}`),
				expDB: []*covidtracker.Risk{{
					ID:              "1",
					RiskLevel:       .5,
					ConfidenceLevel: .5,
					BySegments: []covidtracker.RiskSegment{
						{
							ID:              covidtracker.RiskSegID("1"),
							Segment:         &covidtracker.Segment{Origin: &mock.Paris, Destination: &mock.Bordeaux},
							ConfidenceLevel: .5,
							RiskLevel:       .5,
						},
					},
					Report: covidtracker.Report{
						Advices: []covidtracker.Statement{{Value: "a wise advice"}},
						Pluses:  []covidtracker.Statement{{Value: "it's pretty good"}},
						Minuses: []covidtracker.Statement{{Value: "it's not good"}},
					},
				}},
			},
		}

		for name, tcase := range tcases {
			t.Logf("case %s... :", name)

			db = []*covidtracker.Risk{}
			risk.Reset()

			got, err := client.Do(tcase.tpl, map[string]interface{}{
				"segs":  tcase.segs,
				"prots": []covidtracker.Protection{},
				"hots":  tcase.hots,
			})

			expResult := &gqlResp{}
			if err := json.Unmarshal(tcase.expRaw, expResult); err != nil {
				t.Fatal(err)
			}

			test.Compare(t, err, nil, name+": unexpected error")
			test.Compare(t, risk.ComputeRiskInvoked, true, name+": estimate invokation is expected")
			test.Compare(t, risk.InsertInvoked, true, name+": insert invokation expected")
			test.Compare(t, got, expResult, name+": unexpected result")
			test.Compare(t, db, tcase.expDB, name+": unexpected inserted results")
		}

	})

}
