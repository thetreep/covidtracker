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

import "time"

type Segment struct {
	ID             SegID          `bson:"_id" json:"-"`
	Origin         *Geo           `bson:"origin" json:"origin"`
	Destination    *Geo           `bson:"destination" json:"destination"`
	Departure      time.Time      `bson:"departure" json:"departure"`
	Arrival        time.Time      `bson:"arrival" json:"arrival"`
	Transportation Transportation `bson:"transportation" json:"transportation"`
	HotelID        *string        `bson:"hotelID" json:"-"`
}

type SegID string
