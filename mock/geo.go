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

package mock

import (
	"github.com/thetreep/covidtracker"
	"github.com/thetreep/toolbox/convert"
)

var (
	strP  = convert.StrP
	intP  = convert.IntP
	Paris = covidtracker.Geo{
		Properties: covidtracker.Properties{
			Name:        strP("Paris"),
			PostalCode:  strP("75056"),
			DepCode:     strP("75"),
			RegionCode:  strP("11"),
			Population:  intP(2190327),
			PostalCodes: []string{"75001", "75002", "75003", "75004", "75005"},
		},
		Type: "Feature",
		Geometry: covidtracker.Geometry{
			Coordinates: []float64{2.3488, 48.8534},
			Type:        "Point",
		},
	}
	Bordeaux = covidtracker.Geo{
		Properties: covidtracker.Properties{
			Name:        strP("Bordeaux"),
			PostalCode:  strP("33063"),
			DepCode:     strP("33"),
			RegionCode:  strP("75"),
			Population:  intP(252040),
			PostalCodes: []string{"33000", "33300", "33100", "33090", "33200"},
		},
		Type: "Feature",
		Geometry: covidtracker.Geometry{
			Coordinates: []float64{-0.5874, 44.8572},
			Type:        "Point",
		},
	}
	ParisStd = covidtracker.Geo{
		// REQUIRED. As per GeoJSON spec.
		Properties: covidtracker.Properties{
			// REQUIRED. Namespace.
			GeoCoding: &covidtracker.GeoCoding{
				// REQUIRED. One of "house", "street", "locality", "city", "region", "country".
				// TODO: make a clean list of common cases, plus make clear that the list
				// isn't meant to be closed.
				Type: "house",
				// OPTIONAL. Result accuracy, in meters.
				Accuracy: 20,
				// RECOMMENDED. Suggested label for the result.
				Label: "My Shoes Shop, 64 rue de Metz 75015 Paris",
				// OPTIONAL. Name of the place.
				Name: "My Shoes Shop",
				// OPTIONAL. Housenumber of the place.
				// TODO: what about the suffix (64A, 64 bis, etc.)?
				HouseNumber: "64",
				// OPTIONAL. Street of the place.
				Street: "Rue de Metz",
				// OPTIONAL. Locality of the place.
				Locality: "Paris",
				// OPTIONAL. Postcode of the place.
				PostCode: "75015",
				// OPTIONAL. City of the place.
				City: "Paris 15",
				// OPTIONAL. District of the place.
				District: nil,
				// OPTIONAL. County of the place.
				County: nil,
				// OPTIONAL. State of the place.
				State: nil,
				// OPTIONAL. Country of the place.
				Country: "France",
				// OPTIONAL. Administratives boundaries the feature is included in,
				// as defined in http://wiki.osm.org/wiki/Key:admin_level#admin_level
				// TODO is there some still generic but less OSMish scheme?
				Admin: covidtracker.AdminLevels{
					Level2: "France",
					Level4: "Ile-De-France",
					Level6: "Paris",
				},
				// OPTIONAL. Geohash encoding of coordinates (see http://geohash.org/site/tips.html).
				Geohash: "Ehugh5oofiToh9aWe3heemu7ighee8",
			},
		},
		// REQUIRED. As per GeoJSON spec.
		Type: "Feature",
		// REQUIRED. As per GeoJSON spec.
		Geometry: covidtracker.Geometry{
			Coordinates: []float64{2.3488, 48.8534},
			Type:        "Point",
		},
	}
	BordeauxStd = covidtracker.Geo{
		// REQUIRED. As per GeoJSON spec.
		Properties: covidtracker.Properties{
			// REQUIRED. Namespace.
			GeoCoding: &covidtracker.GeoCoding{
				// REQUIRED. One of "house", "street", "locality", "city", "region", "country".
				// TODO: make a clean list of common cases, plus make clear that the list
				// isn't meant to be closed.
				Type: "house",
				// OPTIONAL. Result accuracy, in meters.
				Accuracy: 20,
				// RECOMMENDED. Suggested label for the result.
				Label: "My Sesam Shop, 14 rue des bananier 33008 Bordeaux",
				// OPTIONAL. Name of the place.
				Name: "My Sesam Shop",
				// OPTIONAL. Housenumber of the place.
				// TODO: what about the suffix (64A, 64 bis, etc.)?
				HouseNumber: "14",
				// OPTIONAL. Street of the place.
				Street: "Rue des bananiers",
				// OPTIONAL. Locality of the place.
				Locality: "Bordeaux",
				// OPTIONAL. Postcode of the place.
				PostCode: "33000",
				// OPTIONAL. City of the place.
				City: "Bordeaux",
				// OPTIONAL. District of the place.
				District: nil,
				// OPTIONAL. County of the place.
				County: nil,
				// OPTIONAL. State of the place.
				State: nil,
				// OPTIONAL. Country of the place.
				Country: "France",
				// OPTIONAL. Administratives boundaries the feature is included in,
				// as defined in http://wiki.osm.org/wiki/Key:admin_level#admin_level
				// TODO is there some still generic but less OSMish scheme?
				Admin: covidtracker.AdminLevels{
					Level2: "France",
					Level4: "Nouvelle-Aquitaine",
					Level6: "Gironde",
				},
				// OPTIONAL. Geohash encoding of coordinates (see http://geohash.org/site/tips.html).
				Geohash: "Ehugh5oofiToh9aWe3heemu7ighee9",
			},
		},
		// REQUIRED. As per GeoJSON spec.
		Type: "Feature",
		// REQUIRED. As per GeoJSON spec.
		Geometry: covidtracker.Geometry{
			Coordinates: []float64{-0.5667, 44.8333},
			Type:        "Point",
		},
	}
)
