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

package covidtracker

import (
	"time"
)

//Risk is the definition of risk and confidence level of a trip
type Risk struct {
	ID              RiskID        `bson:"_id" json:"id"`
	NoticeDate      time.Time     `bson:"noticeDate" json:"noticeDate"`
	ConfidenceLevel float64       `bson:"confidenceLevel" json:"confidenceLevel"`
	RiskLevel       float64       `bson:"riskLevel" json:"riskLevel"`
	DisplayedRisk   float64       `bson:"displayedRisk" json:"displayedRisk"`
	BySegments      []RiskSegment `bson:"bySegments" json:"bySegments"`
	Report          Report        `bson:"report" json:"report"`
}

type Report struct {
	Minuses []Statement `bson:"minuses" json:"minuses"`
	Pluses  []Statement `bson:"pluses" json:"pluses"`
	Advices []Statement `bson:"advices" json:"advices"`
}

type Statement struct {
	Value    string `bson:"value" json:"value"`
	Category string `bson:"category" json:"category"`
}

//RiskID identifies a Risk
type RiskID string

//RiskSegment is the risk and the confidence level for a given segment
type RiskSegment struct {
	ID RiskSegID `bson:"_id" json:"id"`

	*Segment `bson:"segment" json:"segment"`

	RiskLevel       float64 `bson:"riskLevel" json:"riskLevel"`
	ConfidenceLevel float64 `bson:"confidenceLevel" json:"confidenceLevel"`
	Report          Report  `bson:"report" json:"report"`
}

type RiskParameters struct {
	// Use to splecify that these are the default parameters
	IsDefault bool `bson:"default" json:"default"`

	// The protection factor of sewn mask
	SewnMaskProtect float64 `bson:"sewnMaskProtect" json:"sewnMaskProtect"`

	// The protection factor of surgical mask
	SurgicalMaskProtect float64 `bson:"surgicalMaskProtect" json:"surgicalMaskProtect"`

	// The protection factor of ffpx mask
	FFPXMaskProtect float64 `bson:"ffpxMaskProtect" json:"ffpxMaskProtect"`

	// The protection factor of hydro alcoholic gel
	HydroAlcoholicGelProtect float64 `bson:"hydroAlcoholicGelProtect" json:"hydroAlcoholicGelProtect"`

	// The parameters associated with a scope
	Parameters []*RiskParameter `bson:"parameters" json:"parameters"`
}

func (r *RiskParameters) ByScope() map[ParameterScope]*RiskParameter {
	res := make(map[ParameterScope]*RiskParameter)
	for _, p := range r.Parameters {
		res[p.Scope] = p
	}
	return res
}

type RiskParameter struct {
	// The scope of this risk parameter
	Scope ParameterScope `bson:"scope" json:"scope"`

	// The number of persons with direct projection possible
	NbDirect int `bson:"nbDirect" json:"nbDirect"`

	// The probability of contagion via direct projection with an infectious person
	ProbaContagionDirect float64 `bson:"probaContagionDirect" json:"probaContagionDirect"`

	// The protection factor of mask against direct contagion
	MaskProtectDirect float64 `bson:"maskProtectDirect" json:"maskProtectDirect"`

	// The number of persons with direct contact with the person
	NbContact int `bson:"nbContact" json:"nbContact"`

	// The probability of contagion via direct contact with an infectious person
	ProbaContagionContact float64 `bson:"probaContagionContact" json:"probaContagionContact"`

	// The protection factor of mask against contact contagion
	MaskProtectContact float64 `bson:"maskProtectContact" json:"maskProtectContact"`

	// The protection factor of gel against contact contagion
	GelProtectContact float64 `bson:"gelProtectContact" json:"gelProtectContact"`

	// The number of persons with indirect contact
	NbIndirect int `bson:"nbIndirect" json:"nbIndirect"`

	// The probability of contagion via indirect contact with an infectious person
	ProbaContagionIndirect float64 `bson:"probaContagionIndirect" json:"probaContagionIndirect"`

	// The protection factor of mask against indirect contact contagion
	MaskProtectIndirect float64 `bson:"maskProtectIndirect" json:"maskProtectIndirect"`

	// The protection factor of gel against indirect contact contagion
	GelProtectIndirect float64 `bson:"gelProtectIndirect" json:"gelProtectIndirect"`

	// The Pluses of this kind of segment
	Pluses []string `bson:"pluses" json:"pluses"`

	// The Minuses of this kind of segment
	Minuses []string `bson:"minuses" json:"minuses"`

	// The Advices of this kind of segment
	Advices []string `bson:"advices" json:"advices"`
}

type ParameterScope struct {
	// Transportation optionally represents the transportation of this scope (if not a place)
	Transportation Transportation `bson:"transportation" json:"transportation"`

	// Place optionally represents the place of this scope (if not a transportation)
	Place Place `bson:"place" json:"place"`

	Duration TransportationDuration `bson:"duration" json:"duration"`
}

func (s *ParameterScope) String() string {
	if len(s.Transportation) != 0 {
		return string(s.Transportation)
	}
	if len(s.Place) != 0 {
		return string(s.Place)
	}
	return ""
}

//RiskSegID identifies a RiskSegment
type RiskSegID string

//RiskDAL defines the data access layer of risk data
type RiskDAL interface {
	Get(id RiskID) (*Risk, error)
	Insert(r ...*Risk) error
}

//RiskJob defines the job to implements risk data logic
type RiskJob interface {
	ComputeRisk(segs []Segment, protects []Protection) (*Risk, error)
}

//RiskParametersDAL defines the data access layer of risk parameters
type RiskParametersDAL interface {
	GetDefault() (*RiskParameters, error)
	Insert(p *RiskParameters) error
}
