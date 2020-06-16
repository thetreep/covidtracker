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

import "github.com/thetreep/covidtracker"

type Risk struct {
	//Jobs
	ComputeRiskFn      func(segs []covidtracker.Segment, protects []covidtracker.Protection) (*covidtracker.Risk, error)
	ComputeRiskInvoked bool

	//DAL
	GetFn      func(id covidtracker.RiskID) (*covidtracker.Risk, error)
	GetInvoked bool

	InsertFn      func(r ...*covidtracker.Risk) error
	InsertInvoked bool

	//API
	EstimateFn      func(query interface{}) ([]*covidtracker.Risk, error)
	EstimateInvoked bool
}

var (
	_ covidtracker.RiskJob = &Risk{}
	_ covidtracker.RiskDAL = &Risk{}
)

func (r *Risk) Reset() {
	r.InsertInvoked = false
	r.ComputeRiskInvoked = false
	r.GetInvoked = false
	r.EstimateInvoked = false
}
func (r *Risk) ComputeRisk(segs []covidtracker.Segment, protects []covidtracker.Protection) (*covidtracker.Risk, error) {
	r.ComputeRiskInvoked = true
	return r.ComputeRiskFn(segs, protects)
}
func (r *Risk) Get(id covidtracker.RiskID) (*covidtracker.Risk, error) {
	r.GetInvoked = true
	return r.GetFn(id)
}
func (r *Risk) Insert(risks ...*covidtracker.Risk) error {
	r.InsertInvoked = true
	return r.InsertFn(risks...)
}
func (r *Risk) Estimate(query interface{}) ([]*covidtracker.Risk, error) {
	r.EstimateInvoked = true
	return r.EstimateFn(query)
}
