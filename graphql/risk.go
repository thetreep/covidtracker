package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/thetreep/covidtracker"
)

type RiskHandler struct {
	Job covidtracker.RiskJob
	DAL covidtracker.RiskDAL
	API covidtracker.RiskAPI
}

var _ Configurer = &RiskHandler{}

func (h *RiskHandler) Queries() graphql.Fields {
	return graphql.Fields{"risk": h.Get()}
}

func (h *RiskHandler) Mutations() graphql.Fields {
	return graphql.Fields{"risk": nil}
}

func (h *RiskHandler) Get() *graphql.Field {

	//TODO move these definitions outside, in another file ?
	segment := graphql.NewObject(graphql.ObjectConfig{
		Name: "segment",
		Fields: graphql.Fields{
			"origin":         &graphql.Field{Type: graphql.String},
			"destination":    &graphql.Field{Type: graphql.String},
			"datetime":       &graphql.Field{Type: graphql.DateTime},
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
	protect := graphql.NewObject(graphql.ObjectConfig{
		Name: "protection",
		Fields: graphql.Fields{
			"type": &graphql.Field{Type: graphql.String},
			"name": &graphql.Field{Type: graphql.String},
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
				Type: graphql.NewList(segment),
			},
			"protections": &graphql.ArgumentConfig{
				Type: graphql.NewList(protect),
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
				s, ok := i.(covidtracker.Segment)
				if ok {
					segs = append(segs, s)
				}
			}
			for _, i := range protectsI {
				p, ok := i.(covidtracker.Protection)
				if ok {
					protects = append(protects, p)
				}
			}

			r, err := h.Job.ComputeRisk(segs, protects)
			if err != nil {
				return nil, err
			}

			//TODO implements additionnal logic here

			//TODO save in DB

			return r, nil
		},
	}

}
