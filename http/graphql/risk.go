package graphql

import (
	"encoding/json"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"golang.org/x/net/context"

	"github.com/thetreep/covidtracker"
	"github.com/thetreep/covidtracker/logger"
)

type RiskHandler struct {
	Job covidtracker.RiskJob
	DAL covidtracker.RiskDAL
}

var _ Configurer = &RiskHandler{}

func (h *RiskHandler) Queries() graphql.Fields {
	return graphql.Fields{"risk": h.Estimate()}
}

func (h *RiskHandler) Mutations() graphql.Fields {
	return graphql.Fields{"risk": nil}
}

func (h *RiskHandler) Estimate() *graphql.Field {
	logger.Debug(context.Background(), "Estimate")

	segmentIn := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "segmentIn",
			Fields: graphql.InputObjectConfigFieldMap{
				"origin":         &graphql.InputObjectFieldConfig{Type: graphql.String},
				"destination":    &graphql.InputObjectFieldConfig{Type: graphql.String},
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
				"type": &graphql.InputObjectFieldConfig{Type: graphql.String},
				"name": &graphql.InputObjectFieldConfig{Type: graphql.String},
			},
		},
	)

	//TODO move these definitions outside, in another file ?
	segment := graphql.NewObject(graphql.ObjectConfig{
		Name: "segment",
		Fields: graphql.Fields{
			"origin":         &graphql.Field{Type: graphql.String},
			"destination":    &graphql.Field{Type: graphql.String},
			"departure":      &graphql.Field{Type: graphql.DateTime},
			"arrival":        &graphql.Field{Type: graphql.DateTime},
			"transportation": &graphql.Field{Type: graphql.String},
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
	// protect := graphql.NewObject(graphql.ObjectConfig{
	// 	Name: "protection",
	// 	Fields: graphql.Fields{
	// 		"type": &graphql.Field{Type: graphql.String},
	// 		"name": &graphql.Field{Type: graphql.String},
	// 	},
	// })
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
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			segsI := params.Args["segments"].([]interface{})
			protectsI := params.Args["protections"].([]interface{})

			var (
				segs     []covidtracker.Segment
				protects []covidtracker.Protection
			)

			for _, i := range segsI {
				if m, ok := i.(map[string]interface{}); ok {
					var seg covidtracker.Segment
					if err := convert(m, &seg); err != nil {
						return nil, err
					}
					segs = append(segs, seg)
				}
			}

			if len(segs) == 0 {
				return nil, fmt.Errorf("at least one `segment` is mandatory")
			}

			for _, i := range protectsI {
				if m, ok := i.(map[string]interface{}); ok {
					var prot covidtracker.Protection
					if err := convert(m, &prot); err != nil {
						return nil, err
					}
					protects = append(protects, prot)
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
					nil,
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
