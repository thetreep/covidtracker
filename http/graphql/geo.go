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

func GeoIn() *graphql.InputObject {

	adminIn := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "adminIn",
			Fields: graphql.InputObjectConfigFieldMap{
				"level2": &graphql.InputObjectFieldConfig{Type: graphql.String},
				"level4": &graphql.InputObjectFieldConfig{Type: graphql.String},
				"level6": &graphql.InputObjectFieldConfig{Type: graphql.String},
			},
		},
	)

	geoCodingIn := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "geoCodingIn",
			Fields: graphql.InputObjectConfigFieldMap{
				"type":        &graphql.InputObjectFieldConfig{Type: graphql.String},
				"name":        &graphql.InputObjectFieldConfig{Type: graphql.String},
				"accuracy":    &graphql.InputObjectFieldConfig{Type: graphql.Int},
				"label":       &graphql.InputObjectFieldConfig{Type: graphql.String},
				"housenumber": &graphql.InputObjectFieldConfig{Type: graphql.String},
				"street":      &graphql.InputObjectFieldConfig{Type: graphql.String},
				"locality":    &graphql.InputObjectFieldConfig{Type: graphql.String},
				"postcode":    &graphql.InputObjectFieldConfig{Type: graphql.String},
				"city":        &graphql.InputObjectFieldConfig{Type: graphql.String},
				"geoHash":     &graphql.InputObjectFieldConfig{Type: graphql.String},
				"admin":       &graphql.InputObjectFieldConfig{Type: adminIn},
				"county":      &graphql.InputObjectFieldConfig{Type: graphql.String},
				"country":     &graphql.InputObjectFieldConfig{Type: graphql.String},
				"state":       &graphql.InputObjectFieldConfig{Type: graphql.String},
				"geohash":     &graphql.InputObjectFieldConfig{Type: graphql.String},
			},
		},
	)

	geometryIn := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "geometryIn",
			Fields: graphql.InputObjectConfigFieldMap{
				"coordinates": &graphql.InputObjectFieldConfig{Type: graphql.NewList(graphql.Float)},
				"type":        &graphql.InputObjectFieldConfig{Type: graphql.String},
			},
		},
	)

	propIn := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "propIn",
			Fields: graphql.InputObjectConfigFieldMap{
				"geocoding":       &graphql.InputObjectFieldConfig{Type: geoCodingIn},
				"nom":             &graphql.InputObjectFieldConfig{Type: graphql.String},
				"code":            &graphql.InputObjectFieldConfig{Type: graphql.String},
				"codeDepartement": &graphql.InputObjectFieldConfig{Type: graphql.String},
				"codeRegion":      &graphql.InputObjectFieldConfig{Type: graphql.String},
				"codesPostaux":    &graphql.InputObjectFieldConfig{Type: graphql.NewList(graphql.String)},
				"population":      &graphql.InputObjectFieldConfig{Type: graphql.Int},
			},
		},
	)

	return graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "geoIn",
			Fields: graphql.InputObjectConfigFieldMap{
				"properties": &graphql.InputObjectFieldConfig{Type: propIn},
				"type":       &graphql.InputObjectFieldConfig{Type: graphql.String},
				"geometry":   &graphql.InputObjectFieldConfig{Type: geometryIn},
			},
		},
	)

}

func GeoObj() *graphql.Object {

	admin := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "admin",
			Fields: graphql.Fields{
				"level2": &graphql.Field{Type: graphql.String},
				"level4": &graphql.Field{Type: graphql.String},
				"level6": &graphql.Field{Type: graphql.String},
			},
		},
	)

	geoCoding := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "geoCoding",
			Fields: graphql.Fields{
				"type":        &graphql.Field{Type: graphql.String},
				"name":        &graphql.Field{Type: graphql.String},
				"accuracy":    &graphql.Field{Type: graphql.Int},
				"label":       &graphql.Field{Type: graphql.String},
				"housenumber": &graphql.Field{Type: graphql.String},
				"street":      &graphql.Field{Type: graphql.String},
				"locality":    &graphql.Field{Type: graphql.String},
				"postcode":    &graphql.Field{Type: graphql.String},
				"city":        &graphql.Field{Type: graphql.String},
				"geoHash":     &graphql.Field{Type: graphql.String},
				"admin":       &graphql.Field{Type: admin},
				"county":      &graphql.Field{Type: graphql.String},
				"country":     &graphql.Field{Type: graphql.String},
				"state":       &graphql.Field{Type: graphql.String},
				"geohash":     &graphql.Field{Type: graphql.String},
			},
		},
	)

	geometry := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "geometry",
			Fields: graphql.Fields{
				"coordinates": &graphql.Field{Type: graphql.NewList(graphql.Float)},
				"type":        &graphql.Field{Type: graphql.String},
			},
		},
	)

	prop := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "prop",
			Fields: graphql.Fields{
				"geocoding":       &graphql.Field{Type: geoCoding},
				"nom":             &graphql.Field{Type: graphql.String},
				"code":            &graphql.Field{Type: graphql.String},
				"codeDepartement": &graphql.Field{Type: graphql.String},
				"codeRegion":      &graphql.Field{Type: graphql.String},
				"codesPostaux":    &graphql.Field{Type: graphql.NewList(graphql.String)},
				"population":      &graphql.Field{Type: graphql.Int},
			},
		},
	)

	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "geo",
			Fields: graphql.Fields{
				"properties": &graphql.Field{Type: prop},
				"type":       &graphql.Field{Type: graphql.String},
				"geometry":   &graphql.Field{Type: geometry},
			},
		},
	)

}
