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

package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func dateFilter(d time.Time) bson.M {
	start := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
	end := time.Date(d.Year(), d.Month(), d.Day()+1, 0, 0, 0, 0, time.UTC)
	return bson.M{
		"$gte": start,
		"$lt":  end,
	}
}

func dateRangeFilter(start, end time.Time) bson.M {
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	end = time.Date(end.Year(), end.Month(), end.Day()+1, 0, 0, 0, 0, time.UTC)
	return bson.M{
		"$gte": start,
		"$lt":  end,
	}
}
