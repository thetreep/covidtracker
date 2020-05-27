package graphql

import (
	"github.com/graphql-go/graphql"
)

type PingHandler struct{}

var _ Configurer = &PingHandler{}

func (h *PingHandler) Queries() graphql.Fields {
	return graphql.Fields{"ping": h.Ping()}
}

func (h *PingHandler) Mutations() graphql.Fields {
	return graphql.Fields{"ping": nil}
}

func (h *PingHandler) Ping() *graphql.Field {
	return &graphql.Field{
		Type: graphql.String,
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return "ok", nil
		},
	}
}
