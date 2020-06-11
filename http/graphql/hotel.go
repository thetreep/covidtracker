package graphql

import (
	"context"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"

	"github.com/thetreep/covidtracker"
)

type HotelHandler struct {
	Job covidtracker.HotelJob
	DAL covidtracker.HotelDAL
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
			"city": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"prefix": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			city := params.Args["city"].(string)
			name := params.Args["prefix"].(string)

			hotels, err := h.Job.HotelsByPrefix(city, name)
			if err != nil {
				return nil, fmt.Errorf("%v", err)
			}

			var mHotels []*covidtracker.Hotel
			if mHotels, err = h.DAL.Insert(hotels); err != nil {
				return hotels, gqlerrors.NewError(
					"cannot insert to database",
					nil,
					"",
					nil,
					[]int{},
					err,
				)
			}

			return mHotels, nil
		},
	}
}
