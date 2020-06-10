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
