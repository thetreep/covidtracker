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
