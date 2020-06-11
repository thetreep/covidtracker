package covidtracker_test

import (
	"testing"

	"github.com/thetreep/covidtracker"
	"github.com/thetreep/toolbox/test"
)

func TestGeocoding(t *testing.T) {

	t.Run("check postcode", func(t *testing.T) {

		tcases := map[string]struct {
			code string
			exp  bool
		}{
			"123 nok":      {code: "123", exp: false},
			"a123 nok":     {code: "a123", exp: false},
			"123a nok":     {code: "123a", exp: false},
			"78139 ok":     {code: "78139", exp: true},
			"00230 ok":     {code: "00230", exp: true},
			"90001 ok":     {code: "90001", exp: true},
			"12900013 nok": {code: "12900013", exp: false},
			"empty nok":    {code: "", exp: false},
		}

		for n, tcase := range tcases {
			got, want := covidtracker.Geo{
				Properties: covidtracker.Properties{
					GeoCoding: &covidtracker.GeoCoding{PostCode: tcase.code},
				},
			}.Check() == nil, tcase.exp
			test.Compare(t, got, want, n)
		}
	})
}
