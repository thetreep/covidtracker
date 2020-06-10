package graphql

import (
	"context"
	"fmt"

	"github.com/graphql-go/graphql"

	"github.com/thetreep/covidtracker"
)

type HotelHandler struct {
	Job covidtracker.HotelJob
}

var _ Configurer = &HotelHandler{}

func (h *HotelHandler) Queries() graphql.Fields {
	return graphql.Fields{"hotels": h.Search()}
}

func (h *HotelHandler) Mutations() graphql.Fields {
	return graphql.Fields{"hotels": nil}
}

func (h *HotelHandler) Search() *graphql.Field {
	logger.Debug(context.Background(), "Search")

	hotelType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Hotel",
		Fields: graphql.Fields{
			"ID":            &graphql.Field{Type: graphql.String},
			"Name":          &graphql.Field{Type: graphql.String},
			"Address":       &graphql.Field{Type: graphql.String},
			"City":          &graphql.Field{Type: graphql.String},
			"ZipCode":       &graphql.Field{Type: graphql.String},
			"ImageURL":      &graphql.Field{Type: graphql.String},
			"SanitaryInfos": &graphql.Field{Type: graphql.NewList(graphql.String)},
			"SanitaryNote":  &graphql.Field{Type: graphql.Float},
			"SanitaryNorm":  &graphql.Field{Type: graphql.String},
		},
	})

	return &graphql.Field{
		Type: graphql.NewList(hotelType),
		Args: graphql.FieldConfigArgument{
			"prefix": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			name := params.Args["prefix"].(string)

			hotels, err := h.Job.HotelsByPrefix(name)
			if err != nil {
				return nil, fmt.Errorf("%v", err)
			}

			return hotels, nil
		},
	}
}
