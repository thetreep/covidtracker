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

import "context"

type Logfer interface {
	HasErr(ctx context.Context, err error) bool
	HasErrWithFields(ctx context.Context, fields map[string]interface{}, err error) bool
	Debug(ctx context.Context, str string, vars ...interface{})
	DebugWithFields(ctx context.Context, fields map[string]interface{}, str string, vars ...interface{})
	Info(ctx context.Context, str string, vars ...interface{})
	InfoWithFields(ctx context.Context, fields map[string]interface{}, str string, vars ...interface{})
	Warn(ctx context.Context, str string, vars ...interface{})
	WarnWithFields(ctx context.Context, fields map[string]interface{}, str string, vars ...interface{})
	Error(ctx context.Context, str string, vars ...interface{})
	ErrorWithFields(ctx context.Context, fields map[string]interface{}, str string, vars ...interface{})
	Panic(ctx context.Context, str string, vars ...interface{})
	PanicWithFields(ctx context.Context, fields map[string]interface{}, str string, vars ...interface{})
}
