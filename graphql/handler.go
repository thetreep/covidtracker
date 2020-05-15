package graphql

import (
	"context"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"github.com/thetreep/covidtracker/logger"
)

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
	logger.Debug(context.Background, "graphQL schema created")

	return handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	}), nil

}
