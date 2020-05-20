package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/thetreep/covidtracker/graphql"
	"github.com/thetreep/covidtracker/http"
	"github.com/thetreep/covidtracker/job"
	"github.com/thetreep/covidtracker/mongo"
)

func main() {
	//TODO set env variable
	mongoURL := os.Getenv("THETREEP_COVID_TRACKER_MONGO_URL")
	if mongoURL == "" {
		log.Fatal("missing 'THETREEP_COVID_TRACKER_MONGO_URL' env variable")
	}
	// Connect to database.
	mongo := mongo.NewClient(mongoURL)
	err := mongo.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer mongo.Close()

	j := job.NewJob()
	j.RiskDAL = mongo.Risk()

	riskHandler := &graphql.RiskHandler{}
	riskHandler.Job = j.RiskJob
	riskHandler.DAL = mongo.Risk()

	h, err := graphql.NewHandler(riskHandler)
	if err != nil {
		log.Fatal(err)
	}

	// start http server
	s := http.NewServer()
	s.Handler = h
	if err := s.Open(); err != nil {
		log.Fatal(err)
	}

	// We need to shut down gracefully when the user hits Ctrl-C.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
	sig := <-sigc
	switch sig {
	case syscall.SIGKILL, syscall.SIGQUIT:
		// Go for the program exit. Don't wait for the server to finish.
		fmt.Println("Received SIGTERM or SIGQUIT, exiting without waiting for the web server to shut down")
		return
	case syscall.SIGTERM, syscall.SIGINT:
		// Stop the server gracefully.
		fmt.Println("Received SIGINT or SIGTERM, shutting down web server gracefully")
	}
	s.Close()

}
