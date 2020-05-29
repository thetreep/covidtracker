package mock

import "github.com/thetreep/covidtracker"

var (
	Paris = covidtracker.Geo{
		// REQUIRED. As per GeoJSON spec.
		Properties: covidtracker.Properties{
			// REQUIRED. Namespace.
			GeoCoding: covidtracker.GeoCoding{
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
	Bordeaux = covidtracker.Geo{
		// REQUIRED. As per GeoJSON spec.
		Properties: covidtracker.Properties{
			// REQUIRED. Namespace.
			GeoCoding: covidtracker.GeoCoding{
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
