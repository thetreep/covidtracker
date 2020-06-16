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

package graphql

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"golang.org/x/net/context"

	"github.com/thetreep/covidtracker"
)

type RiskHandler struct {
	Job covidtracker.RiskJob
	DAL covidtracker.RiskDAL
}

var _ Configurer = &RiskHandler{}

type HotelInput struct {
	ID       string    `json:"id"`
	Arrival  time.Time `json:"arrival"`
	NbNights int       `json:"nbNights"`
}

func (h *RiskHandler) Queries() graphql.Fields {
	return graphql.Fields{"risk": h.Estimate()}
}

func (h *RiskHandler) Mutations() graphql.Fields {
	return graphql.Fields{"risk": nil}
}

func (h *RiskHandler) Estimate() *graphql.Field {
	logger.Debug(context.Background(), "Estimate")

	geoIn := GeoIn()
	segmentIn := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "segmentIn",
			Fields: graphql.InputObjectConfigFieldMap{
				"origin":         &graphql.InputObjectFieldConfig{Type: geoIn},
				"destination":    &graphql.InputObjectFieldConfig{Type: geoIn},
				"departure":      &graphql.InputObjectFieldConfig{Type: graphql.DateTime},
				"arrival":        &graphql.InputObjectFieldConfig{Type: graphql.DateTime},
				"transportation": &graphql.InputObjectFieldConfig{Type: graphql.String},
			},
		},
	)
	protIn := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "protectionIn",
			Fields: graphql.InputObjectConfigFieldMap{
				"type":     &graphql.InputObjectFieldConfig{Type: graphql.String},
				"quantity": &graphql.InputObjectFieldConfig{Type: graphql.Int},
			},
		},
	)

	hotelIn := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "hotelIn",
			Fields: graphql.InputObjectConfigFieldMap{
				"id":       &graphql.InputObjectFieldConfig{Type: graphql.String},
				"nbNights": &graphql.InputObjectFieldConfig{Type: graphql.Int},
				"arrival":  &graphql.InputObjectFieldConfig{Type: graphql.DateTime},
			},
		},
	)

	geoObj := GeoObj()
	segment := graphql.NewObject(graphql.ObjectConfig{
		Name: "segment",
		Fields: graphql.Fields{
			"origin":         &graphql.Field{Type: geoObj},
			"destination":    &graphql.Field{Type: geoObj},
			"departure":      &graphql.Field{Type: graphql.DateTime},
			"arrival":        &graphql.Field{Type: graphql.DateTime},
			"transportation": &graphql.Field{Type: graphql.String},
			"hotelID":        &graphql.Field{Type: graphql.String},
		},
	})

	segmentRisk := graphql.NewObject(graphql.ObjectConfig{
		Name: "segmentRisk",
		Fields: graphql.Fields{
			"segment":         &graphql.Field{Type: segment},
			"riskLevel":       &graphql.Field{Type: graphql.Float},
			"condidenceLevel": &graphql.Field{Type: graphql.Float},
		},
	})

	statement := graphql.NewObject(graphql.ObjectConfig{
		Name: "statement",
		Fields: graphql.Fields{
			"value":    &graphql.Field{Type: graphql.String},
			"category": &graphql.Field{Type: graphql.String}},
	})
	report := graphql.NewObject(graphql.ObjectConfig{
		Name: "report",
		Fields: graphql.Fields{
			"minuses": &graphql.Field{Type: graphql.NewList(statement)},
			"pluses":  &graphql.Field{Type: graphql.NewList(statement)},
			"advices": &graphql.Field{Type: graphql.NewList(statement)},
		},
	})
	risk := graphql.NewObject(graphql.ObjectConfig{
		Name: "Risk",
		Fields: graphql.Fields{
			"noticeDate":      &graphql.Field{Type: graphql.DateTime},
			"riskLevel":       &graphql.Field{Type: graphql.Float},
			"confidenceLevel": &graphql.Field{Type: graphql.Float},
			"bySegments":      &graphql.Field{Type: graphql.NewList(segmentRisk)},
			"report":          &graphql.Field{Type: report},
		},
	})

	return &graphql.Field{
		Type: risk,
		Args: graphql.FieldConfigArgument{
			"segments": &graphql.ArgumentConfig{
				Type: graphql.NewList(segmentIn),
			},
			"protections": &graphql.ArgumentConfig{
				Type: graphql.NewList(protIn),
			},
			"hotels": &graphql.ArgumentConfig{
				Type: graphql.NewList(hotelIn),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			segsI, segOK := params.Args["segments"].([]interface{})
			if !segOK {
				logger.Warn(context.Background(), "impossible to read segments input %v", params.Args["segments"])
			}
			hotelsI, hotOK := params.Args["hotels"].([]interface{})
			if !hotOK {
				logger.Warn(context.Background(), "impossible to read hotels input %v", params.Args["hotels"])
			}
			protectsI, protOK := params.Args["protections"].([]interface{})
			if !protOK {
				logger.Warn(context.Background(), "impossible to read protections input %v", params.Args["protections"])
			}

			var (
				segs     []covidtracker.Segment
				protects []covidtracker.Protection
			)

			if segOK {
				for _, i := range segsI {
					if m, ok := i.(map[string]interface{}); ok {
						var seg covidtracker.Segment
						if err := convert(m, &seg); err != nil {
							return nil, err
						}
						segs = append(segs, seg)
					}
				}
			}

			if hotOK {
				for _, i := range hotelsI {
					if m, ok := i.(map[string]interface{}); ok {
						var hin HotelInput
						if err := convert(m, &hin); err != nil {
							return nil, err
						}

						segs = append(segs, covidtracker.Segment{
							Departure: hin.Arrival,
							Arrival:   hin.Arrival.AddDate(0, 0, hin.NbNights),
							HotelID:   &hin.ID,
						})
					}
				}
			}

			if len(segs) == 0 {
				return nil, fmt.Errorf("at least one `segment` or `hotel` is mandatory")
			}

			if protOK {
				for _, i := range protectsI {
					if m, ok := i.(map[string]interface{}); ok {
						var prot covidtracker.Protection
						if err := convert(m, &prot); err != nil {
							return nil, err
						}
						protects = append(protects, prot)
					}
				}
			}

			r, err := h.Job.ComputeRisk(segs, protects)
			if err != nil {
				return nil, err
			}

			if err := h.DAL.Insert(r); err != nil {
				return r, gqlerrors.NewError(
					"cannot insert to database",
					nil,
					"",
					nil,
					[]int{},
					err,
				)
			}

			return r, nil
		},
	}

}

func convert(m map[string]interface{}, output interface{}) error {
	str, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(str, &output); err != nil {
		return err
	}
	return nil
}
