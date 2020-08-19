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

package datagouv

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/pkg/errors"
	"github.com/thetreep/covidtracker"
)

type (
	ResourceID string
	DatasetID  string
)

const (
	DatagouvBase = "https://www.data.gouv.fr"

	EmergencyID       ResourceID = "emergency"
	CaseID            ResourceID = "service-with-case"
	HospitalizationID ResourceID = "hospitalization"
	ScreeningID       ResourceID = "screening"
	IndicatorID       ResourceID = "indicator"

	EmergencyDataset DatasetID = "donnees-des-urgences-hospitalieres-et-de-sos-medecins-relatives-a-lepidemie-de-covid-19"
	HospDataset      DatasetID = "donnees-hospitalieres-relatives-a-lepidemie-de-covid-19"
	ScreeningDataset DatasetID = "donnees-relatives-aux-tests-de-depistage-de-covid-19-realises-en-laboratoire-de-ville"
	IndicDataset     DatasetID = "indicateurs-dactivite-epidemique-covid-19-par-departement"
)

var (
	caseRegex      = regexp.MustCompile(`donnees-hospitalieres-etablissements-covid19-\d{4}-\d{2}-\d{2}-\d{2}h\d{2}.csv`)
	emerRegex      = regexp.MustCompile(`sursaud-corona-quot-dep-\d{4}-\d{2}-\d{2}-\d{2}h\d{2}.csv`)
	hospRegex      = regexp.MustCompile(`donnees-hospitalieres-covid19-\d{4}-\d{2}-\d{2}-\d{2}h\d{2}.csv`)
	indicRegex     = regexp.MustCompile(`donnees-carte-synthese-tricolore.csv`)
	screeningRegex = regexp.MustCompile(`donnees-tests-covid19-labo-quotidien-\d{4}-\d{2}-\d{2}-\d{2}h\d{2}.csv`)
)

type Dataset struct {
	Acronym       interface{}   `json:"acronym"`
	Archived      interface{}   `json:"archived"`
	Badges        []interface{} `json:"badges"`
	CreatedAt     string        `json:"created_at"`
	Deleted       interface{}   `json:"deleted"`
	Description   string        `json:"description"`
	Extras        struct{}      `json:"extras"`
	Frequency     string        `json:"frequency"`
	FrequencyDate interface{}   `json:"frequency_date"`
	ID            string        `json:"id"`
	LastModified  string        `json:"last_modified"`
	LastUpdate    string        `json:"last_update"`
	License       string        `json:"license"`
	Metrics       struct {
		Discussions int `json:"discussions"`
		Followers   int `json:"followers"`
		Issues      int `json:"issues"`
		Reuses      int `json:"reuses"`
		Views       int `json:"views"`
	} `json:"metrics"`
	Organization struct {
		Acronym       interface{} `json:"acronym"`
		Class         string      `json:"class"`
		ID            string      `json:"id"`
		Logo          string      `json:"logo"`
		LogoThumbnail string      `json:"logo_thumbnail"`
		Name          string      `json:"name"`
		Page          string      `json:"page"`
		Slug          string      `json:"slug"`
		URI           string      `json:"uri"`
	} `json:"organization"`
	Owner            interface{}   `json:"owner"`
	Page             string        `json:"page"`
	Private          bool          `json:"private"`
	Resources        []Resource    `json:"resources"`
	Slug             string        `json:"slug"`
	Spatial          interface{}   `json:"spatial"`
	Tags             []interface{} `json:"tags"`
	TemporalCoverage interface{}   `json:"temporal_coverage"`
	Title            string        `json:"title"`
	URI              string        `json:"uri"`
}

type Resource struct {
	Checksum struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"checksum"`
	CreatedAt   string      `json:"created_at"`
	Description interface{} `json:"description"`
	Extras      struct {
	} `json:"extras"`
	Filesize     int    `json:"filesize"`
	Filetype     string `json:"filetype"`
	Format       string `json:"format"`
	ID           string `json:"id"`
	LastModified string `json:"last_modified"`
	Latest       string `json:"latest"`
	Metrics      struct {
		Views int `json:"views"`
	} `json:"metrics"`
	Mime       string      `json:"mime"`
	PreviewURL interface{} `json:"preview_url"`
	Published  string      `json:"published"`
	Title      string      `json:"title"`
	Type       string      `json:"type"`
	URL        string      `json:"url"`
}

type Service struct {
	Ctx      context.Context
	BasePath string

	log covidtracker.Logfer
}

func NewService(ctx context.Context, l covidtracker.Logfer) *Service {
	return &Service{Ctx: ctx, BasePath: DatagouvBase, log: l}
}

func (s *Service) GetCSV(id ResourceID) (*csv.Reader, func() error, error) {
	//get dataset resource URL
	res, err := s.getResourceURL(id)
	if err != nil {
		return nil, nil, err
	}

	// download the resource file
	resp, err := http.Get(string(res))
	if err != nil {
		return nil, nil, err
	}

	r, err := s.newReader(resp.Body, id)
	if err != nil {
		return nil, nil, err
	}

	return r, resp.Body.Close, nil
}

func (s *Service) getResourceURL(resource ResourceID) (string, error) {

	var (
		id  DatasetID
		rgx *regexp.Regexp
	)

	switch resource {
	case EmergencyID:
		id = EmergencyDataset
		rgx = emerRegex
	case CaseID:
		id = HospDataset
		rgx = caseRegex
	case HospitalizationID:
		id = HospDataset
		rgx = hospRegex
	case ScreeningID:
		id = ScreeningDataset
		rgx = screeningRegex
	case IndicatorID:
		id = IndicDataset
		rgx = indicRegex
	}

	res, err := http.Get(s.BasePath + "/api/1/datasets/" + string(id))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	//format response
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, res.Body); err != nil {
		return "", errors.Wrap(err, "reading body")
	}
	dataset := &Dataset{}
	if err := json.NewDecoder(&buf).Decode(&dataset); err != nil {
		return "", errors.Wrap(err, "decoding response")
	}

	// try to find the corresponding resource
	for _, r := range dataset.Resources {
		if rgx.MatchString(r.Title) {
			return r.Latest, nil
		}
	}

	return "", fmt.Errorf("no resource %q found (dataset=%q)", resource, id)
}

func (s *Service) newReader(r io.Reader, res ResourceID) (*csv.Reader, error) {

	reader := csv.NewReader(r)
	switch res {
	case IndicatorID:
		reader.Comma = ','
	case CaseID, HospitalizationID, ScreeningID, EmergencyID:
		reader.Comma = ';'
	default:
		return nil, fmt.Errorf("unsupported dataset %s", res)
	}

	return reader, nil
}

func (s *Service) handleParsingErr(err error, name, col string) error {
	if err != nil {
		s.log.Warn(s.Ctx, "%s : cannot parse %q column (%v)\n", name, col, err)
	}
	return err
}
