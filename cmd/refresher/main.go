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

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/robfig/cron"
	"github.com/thetreep/covidtracker/job"
	"github.com/thetreep/covidtracker/job/datagouv"
	"github.com/thetreep/covidtracker/logger"
	"github.com/thetreep/covidtracker/mongo"
)

func main() {

	//TODO set env variable
	mongoURL := os.Getenv("THETREEP_COVIDTRACKER_MONGO_URL")
	if mongoURL == "" {
		log.Fatal("missing 'THETREEP_COVIDTRACKER_MONGO_URL' env variable")
	}
	// Connect to database.
	mongo := mongo.NewClient(mongoURL)
	err := mongo.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer mongo.Close()

	log := &logger.Logger{}

	j := job.NewJob(log, datagouv.NewService(context.Background(), log))
	j.RiskDAL = mongo.Risk()

	c := cron.New()
	c.AddFunc("@midnight", func() {
		fmt.Print(j.RefreshJob.Refresh(mongo.Case(), mongo.Emergency(), mongo.Hospitalization(), mongo.Indicator(), mongo.Screening()))
	})

	// launch it once to update data at first launch
	err = j.RefreshJob.Refresh(mongo.Case(), mongo.Emergency(), mongo.Hospitalization(), mongo.Indicator(), mongo.Screening())
	log.HasErr(context.Background(), err)

	c.Run()

}
