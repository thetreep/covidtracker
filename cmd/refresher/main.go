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

	j := job.NewJob(datagouv.NewService(context.Background(), &logger.Logger{}))
	j.RiskDAL = mongo.Risk()

	c := cron.New()
	c.AddFunc("@midnight", func() {
		fmt.Print(j.RefreshJob.Refresh(mongo.Case(), mongo.Emergency(), mongo.Hospitalization(), mongo.Indicator(), mongo.Screening()))
	})

}
