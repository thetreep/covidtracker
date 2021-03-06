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

package graphql_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/thetreep/covidtracker"
	"github.com/thetreep/covidtracker/http/graphql"
	"github.com/thetreep/covidtracker/mock"

	"github.com/thetreep/toolbox/test"
)

func TestSearch(t *testing.T) {
	ctx := context.Background()

	hotel := mock.Hotel{}
	h, err := graphql.NewHandler(&graphql.HotelHandler{Job: &hotel, DAL: &hotel})
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(h)
	defer server.Close()

	client := NewClient(ctx, server.URL)

	allHotel := `
	query hotels($city: String!, $prefix: String!) {
		hotels(
			city:$city
			prefix:$prefix,
		) {
			Name, Address, City, ZipCode, ImageURL, SanitaryInfos, SanitaryNote, SanitaryNorm
		}
	}`

	nameHotel := `
	query hotels($city: String!, $prefix: String!) {
		hotels(
			city:$city
			prefix:$prefix,
		) {
			Name
		}
	}`

	noteHotel := `
	query hotels($city: String!, $prefix: String!) {
		hotels(
			city:$city
			prefix:$prefix,
		) {
			SanitaryNote
		}
	}`

	t.Run("hotels error", func(t *testing.T) {
		hotel.HotelsByPrefixFn = func(city string, prefix string) ([]*covidtracker.Hotel, error) {
			return nil, fmt.Errorf("search hotels error")
		}

		got, err := client.Do(allHotel, map[string]interface{}{
			"city":   "",
			"prefix": "",
		})
		expected := &gqlResp{
			Data:   map[string]interface{}{"hotels": nil},
			Errors: []gqlErr{{Message: "search hotels error"}},
		}
		test.Compare(t, err.Error(), "search hotels error", "unexpected error")
		test.Compare(t, got, expected, "unexpected result")
		test.Compare(t, hotel.HotelsByPrefixInvoked, true, "HotelsByPrefix invokation is expected")
	})

	t.Run("Hotels found", func(t *testing.T) {

		hotel.HotelsByPrefixFn = func(city string, prefix string) ([]*covidtracker.Hotel, error) {
			h := []*covidtracker.Hotel{
				&covidtracker.Hotel{
					Address:  "69, Boulevard Sakakini",
					City:     "Marseille",
					ImageURL: "https://bookings.cdsgroupe.com/photos/Search/FR/ACC/251/ACC2514.jpg",
					Name:     "Ibis Budget Marseille Timone",
					SanitaryInfos: []string{
						"Enregistrement & Règlement en ligne",
						"Distanciation sociale & sens de circulation",
						"Formation des équipes internes aux mesures internes",
						"Horaires & Messages de nettoyages des chambres",
						"Renforcement du nettoyage du linge",
						"Procédures de nettoyage renforcées des points de contacts en chambre",
						"Port du masque par le personnel",
						"Mise à disposition de gel hydro-alcoolique",
						"Nettoyage renforcé des lieux de passage",
					},
					SanitaryNorm: "Accor - All Safe",
					SanitaryNote: 7,
					ZipCode:      "13005",
				},
			}
			return h, nil
		}

		db := []*covidtracker.Hotel{}
		hotel.InsertFn = func(hotels []*covidtracker.Hotel) ([]*covidtracker.Hotel, error) {
			for i, hotel := range hotels {
				hotel.ID = covidtracker.HotelID(fmt.Sprint(i + 1))
				db = append(db, hotel)
			}
			return hotels, nil
		}

		tcases := map[string]struct {
			tpl    string
			expRaw []byte
		}{
			"full query": {
				tpl: allHotel,
				expRaw: []byte(`{
					"data": {
						"hotels": [
							{
								"Address": "69, Boulevard Sakakini",
								"City": "Marseille",
								"ImageURL": "https://bookings.cdsgroupe.com/photos/Search/FR/ACC/251/ACC2514.jpg",
								"Name": "Ibis Budget Marseille Timone",
								"SanitaryInfos": [
									"Enregistrement & Règlement en ligne",
									"Distanciation sociale & sens de circulation",
									"Formation des équipes internes aux mesures internes",
									"Horaires & Messages de nettoyages des chambres",
									"Renforcement du nettoyage du linge",
									"Procédures de nettoyage renforcées des points de contacts en chambre",
									"Port du masque par le personnel",
									"Mise à disposition de gel hydro-alcoolique",
									"Nettoyage renforcé des lieux de passage"
								],
								"SanitaryNorm": "Accor - All Safe",
								"SanitaryNote": 7,
								"ZipCode": "13005"
							}
						]
					}
				}`),
			},
			"only name": {
				tpl: nameHotel,
				expRaw: []byte(`{
					"data": {
						"hotels": [
							{
								"Name": "Ibis Budget Marseille Timone"
							}
						]
					}
				}`),
			},
			"only sanitary note": {
				tpl: noteHotel,
				expRaw: []byte(`{
					"data": {
						"hotels": [
							{
								"SanitaryNote": 7
							}
						]
					}
				}`),
			},
		}

		expDB := []*covidtracker.Hotel{{
			ID:       "1",
			Address:  "69, Boulevard Sakakini",
			City:     "Marseille",
			ImageURL: "https://bookings.cdsgroupe.com/photos/Search/FR/ACC/251/ACC2514.jpg",
			Name:     "Ibis Budget Marseille Timone",
			SanitaryInfos: []string{
				"Enregistrement & Règlement en ligne",
				"Distanciation sociale & sens de circulation",
				"Formation des équipes internes aux mesures internes",
				"Horaires & Messages de nettoyages des chambres",
				"Renforcement du nettoyage du linge",
				"Procédures de nettoyage renforcées des points de contacts en chambre",
				"Port du masque par le personnel",
				"Mise à disposition de gel hydro-alcoolique",
				"Nettoyage renforcé des lieux de passage",
			},
			SanitaryNorm: "Accor - All Safe",
			SanitaryNote: 7,
			ZipCode:      "13005",
		}}

		for name, tcase := range tcases {
			t.Logf("case %s... :", name)

			db = []*covidtracker.Hotel{}

			hotel.Reset()

			got, err := client.Do(tcase.tpl, map[string]interface{}{
				"city":   "Marseille",
				"prefix": "Ibis Budget",
			})

			expResult := &gqlResp{}
			if err := json.Unmarshal(tcase.expRaw, expResult); err != nil {
				t.Fatal("test", err)
			}

			test.Compare(t, err, nil, name+": unexpected error")
			test.Compare(t, hotel.HotelsByPrefixInvoked, true, name+": hotelsByPrefix invokation is expected")
			test.Compare(t, hotel.InsertInvoked, true, name+": insert invokation expected")
			test.Compare(t, db, expDB, name+": unexpected inserted results")
			test.Compare(t, got, expResult, name+": unexpected result")
		}
	})
}
