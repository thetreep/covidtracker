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
	return graphql.Fields{}
}

func (h *RiskHandler) Get() *graphql.Field {

	risk := graphql.NewObject(graphql.ObjectConfig{
		Name: "Risk",
		Fields: graphql.Fields{
			"noticeDate":      &graphql.Field{Type: graphql.DateTime},
			"riskLevel":       &graphql.Field{Type: graphql.Float},
			"confidenceLevel": &graphql.Field{Type: graphql.Float},
			//TODO add routes and protections ?
		},
	})

	return &graphql.Field{
		Type: risk,
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {

			r, err := h.Job.ComputeRisk()
			if err != nil {
				return nil, err
			}

			//TODO implements additionnal logic here

			//TODO save in DB

			return r, nil
		},
	}

}
