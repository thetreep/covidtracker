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
	"context"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/handler"

	log "github.com/thetreep/covidtracker/logger"
)

var logger = log.Logger{}

type Configurer interface {
	Queries() graphql.Fields
	Mutations() graphql.Fields
}

func NewHandler(configs ...Configurer) (http.Handler, error) {

	queries := make(graphql.Fields)
	mutations := make(graphql.Fields)

	for _, config := range configs {
		for key, query := range config.Queries() {
			queries[key] = query
		}
		for key, mutation := range config.Mutations() {
			mutations[key] = mutation
		}
	}

	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name:   "RootQuery",
			Fields: queries,
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name:   "RootMutation",
			Fields: mutations,
		}),
	}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		return nil, err
	}

	logger.Debug(context.Background(), "graphQL schema created")

	return handler.New(&handler.Config{
		FormatErrorFn: func(err error) gqlerrors.FormattedError {
			logger.Info(context.Background(), "error occured: %s", err)
			return gqlerrors.FormatError(err)
		},
		Schema: &schema,
		Pretty: true,
	}), nil

}
