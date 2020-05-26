package graphql_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/thetreep/covidtracker"
	"github.com/thetreep/toolbox/test"

	"github.com/thetreep/covidtracker/http/graphql"
	"github.com/thetreep/covidtracker/mock"
)

func TestEstimate(t *testing.T) {
	ctx := context.Background()

	risk := mock.Risk{}
	h, err := graphql.NewHandler(&graphql.RiskHandler{DAL: &risk, Job: &risk})
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(h)
	defer server.Close()

	client := NewClient(ctx, server.URL)

	template := `
	query ($segs: [segmentIn], $prots: [protectionIn]) {
		risk(
			segments:$segs,
			protections:$prots
		) {
			bySegments {
				riskLevel
				segment {
					origin
					destination
				}
			}
			riskLevel
			confidenceLevel
			report {
				minuses {value}
				pluses {value}
				advices {value}
			}
		}
	}`

	t.Run("no segment provided", func(t *testing.T) {
		got, err := client.Do(template, map[string]interface{}{
			"segs":  []covidtracker.Segment{},
			"prots": []covidtracker.Protection{},
		})
		expected := &gqlResp{
			Data:   map[string]interface{}{"risk": nil},
			Errors: []gqlErr{{Message: "at least one `segment` is mandatory"}},
		}
		test.Compare(t, err.Error(), "at least one `segment` is mandatory", "unexpected error")
		test.Compare(t, got, expected, "unexpected result")
		test.Compare(t, risk.ComputeRiskInvoked, false, "estimate invokation unexpected")
		test.Compare(t, risk.InsertInvoked, false, "insert invokation unexpected")
	})

	t.Run("computation error", func(t *testing.T) {
		risk.ComputeRiskFn = func(segs []covidtracker.Segment, protects []covidtracker.Protection) (*covidtracker.Risk, error) {
			return nil, fmt.Errorf("computation error")
		}

		got, err := client.Do(template, map[string]interface{}{
			"segs":  []covidtracker.Segment{{Origin: "paris", Destination: "lyon"}},
			"prots": []covidtracker.Protection{},
		})
		expected := &gqlResp{
			Data:   map[string]interface{}{"risk": nil},
			Errors: []gqlErr{{Message: "computation error"}},
		}
		test.Compare(t, err.Error(), "computation error", "unexpected error")
		test.Compare(t, got, expected, "unexpected result")
		test.Compare(t, risk.ComputeRiskInvoked, true, "estimate invokation is expected")
		test.Compare(t, risk.InsertInvoked, false, "insert invokation unexpected")
	})

}
