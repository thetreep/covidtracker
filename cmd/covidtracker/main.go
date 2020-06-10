package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/thetreep/covidtracker"
	"github.com/thetreep/covidtracker/http"
	"github.com/thetreep/covidtracker/http/graphql"
	"github.com/thetreep/covidtracker/job"
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

	j := job.NewJob()
	j.RiskDAL = mongo.Risk()
	j.ParametersDAL = mongo.Parameters()

	pingHandler := &graphql.PingHandler{}

	riskHandler := &graphql.RiskHandler{}
	riskHandler.Job = j.Risk()
	riskHandler.DAL = mongo.Risk()

	if err := createDefaultParametersIfMissing(mongo.Parameters()); err != nil {
		log.Fatal(err)
	}

	gql, err := graphql.NewHandler(pingHandler, riskHandler)
	if err != nil {
		log.Fatal(err)
	}

	// start http server
	s := http.NewServer()
	s.AddHandler(gql, "/graphql")
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

func createDefaultParametersIfMissing(dal covidtracker.RiskParametersDAL) error {
	_, err := dal.GetDefault()
	if err == nil {
		return nil // default parameters are already existing
	}
	if err != covidtracker.ErrNoParametersDefined {
		return err
	}
	defaultParams := &covidtracker.RiskParameters{
		IsDefault: true,
		ParametersByScope: map[covidtracker.ParameterScope]covidtracker.RiskParameter{
			{Transportation: covidtracker.Aircraft, Duration: covidtracker.Long}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{"Segment avion long"},
				Advices:                []string{},
			},
			{Transportation: covidtracker.Aircraft, Duration: covidtracker.Normal}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{},
				Advices:                []string{},
			},
			{Transportation: covidtracker.Aircraft, Duration: covidtracker.Short}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{"Segment avion court"},
				Minuses:                []string{},
				Advices:                []string{},
			},
			{Transportation: covidtracker.TGV, Duration: covidtracker.Long}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{"Segment train long"},
				Advices:                []string{},
			},
			{Transportation: covidtracker.TGV, Duration: covidtracker.Normal}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{},
				Advices:                []string{},
			},
			{Transportation: covidtracker.TGV, Duration: covidtracker.Short}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{"Segment train court"},
				Minuses:                []string{},
				Advices:                []string{},
			},
			{Transportation: covidtracker.TER, Duration: covidtracker.Long}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{"Segment train long"},
				Advices:                []string{},
			},
			{Transportation: covidtracker.TER, Duration: covidtracker.Normal}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{},
				Advices:                []string{},
			},
			{Transportation: covidtracker.TER, Duration: covidtracker.Short}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{"Segment train court"},
				Minuses:                []string{},
				Advices:                []string{},
			},
			{Transportation: covidtracker.CarSolo, Duration: covidtracker.Long}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{"Vous êtes seul(e) dans la voiture"},
				Minuses:                []string{"Segment voiture long, il faudra probablement s'arrêter à une pompe à essence"},
				Advices:                []string{"Lavez vous bien les mains si vous prenez de l'essence"},
			},
			{Transportation: covidtracker.CarSolo, Duration: covidtracker.Normal}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{"Vous êtes seul(e) dans la voiture"},
				Minuses:                []string{},
				Advices:                []string{"Lavez vous bien les mains si vous prenez de l'essence"},
			},
			{Transportation: covidtracker.CarSolo, Duration: covidtracker.Short}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{"Segment voiture court", "Vous êtes seul(e) dans la voiture"},
				Minuses:                []string{},
				Advices:                []string{},
			},
			{Transportation: covidtracker.CarDuo, Duration: covidtracker.Long}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{"Segment voiture long, il faudra probablement s'arrêter à une pompe à essence", "Vous êtes plusieurs dans voiture"},
				Advices:                []string{"Lavez vous bien les mains si vous prenez de l'essence"},
			},
			{Transportation: covidtracker.CarDuo, Duration: covidtracker.Normal}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{"Vous êtes plusieurs dans voiture"},
				Advices:                []string{"Lavez vous bien les mains si vous prenez de l'essence"},
			},
			{Transportation: covidtracker.CarDuo, Duration: covidtracker.Short}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{"Segment voiture court"},
				Minuses:                []string{"Vous êtes plusieurs dans voiture"},
				Advices:                []string{},
			},
			{Transportation: covidtracker.CarGroup, Duration: covidtracker.Long}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{"Segment voiture long, il faudra probablement s'arrêter à une pompe à essence", "Vous êtes plusieurs dans voiture"},
				Advices:                []string{"Lavez vous bien les mains si vous prenez de l'essence"},
			},
			{Transportation: covidtracker.CarGroup, Duration: covidtracker.Normal}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{"Vous êtes plusieurs dans voiture"},
				Advices:                []string{"Lavez vous bien les mains si vous prenez de l'essence"},
			},
			{Transportation: covidtracker.CarGroup, Duration: covidtracker.Short}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{"Segment voiture court"},
				Minuses:                []string{"Vous êtes plusieurs dans voiture"},
				Advices:                []string{},
			},
			{Transportation: covidtracker.PublicTransports, Duration: covidtracker.Long}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{},
				Advices:                []string{},
			},
			{Transportation: covidtracker.PublicTransports, Duration: covidtracker.Normal}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{},
				Advices:                []string{},
			},
			{Transportation: covidtracker.PublicTransports, Duration: covidtracker.Short}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{},
				Advices:                []string{},
			},
			{Transportation: covidtracker.TaxiSolo, Duration: covidtracker.Normal}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{},
				Advices:                []string{},
			},
			{Transportation: covidtracker.TaxiGroup, Duration: covidtracker.Normal}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{"Vous êtes plusieurs passagers dans le taxi"},
				Advices:                []string{},
			},
			{Transportation: covidtracker.Scooter, Duration: covidtracker.Normal}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{},
				Advices:                []string{},
			},
			{Transportation: covidtracker.Bike, Duration: covidtracker.Normal}: {
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{},
				Advices:                []string{},
			},
		},
	}
	return dal.Insert(defaultParams)
}
