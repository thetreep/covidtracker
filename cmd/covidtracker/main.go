package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/thetreep/covidtracker"
	"github.com/thetreep/covidtracker/http"
	"github.com/thetreep/covidtracker/http/graphql"
	"github.com/thetreep/covidtracker/job"
	"github.com/thetreep/covidtracker/job/cds"
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
	ctx := context.Background()
	j := job.NewJob(log, datagouv.NewService(context.Background(), log))
	j.RiskDAL = mongo.Risk()
	j.RiskParametersDAL = mongo.RiskParameters()
	j.EmergencyDAL = mongo.Emergency()
	j.HotelDAL = mongo.Hotel()

	pingHandler := &graphql.PingHandler{}

	riskHandler := &graphql.RiskHandler{}
	riskHandler.Job = j.Risk()
	riskHandler.DAL = mongo.Risk()

	if err := createDefaultParametersIfMissing(mongo.RiskParameters()); err != nil {
		log.Fatal(ctx, err.Error())
	}

	cds.Init()
	HotelHandler := &graphql.HotelHandler{}
	HotelHandler.Job = j.Hotels()
	HotelHandler.DAL = mongo.Hotel()

	gql, err := graphql.NewHandler(pingHandler, riskHandler, HotelHandler)
	if err != nil {
		log.Fatal(ctx, err.Error())
	}

	// start http server
	s := http.NewServer()
	s.AddHandler(gql, "/graphql")
	if err := s.Open(); err != nil {
		log.Fatal(ctx, err.Error())
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
		IsDefault:                true,
		SewnMaskProtect:          0.71,
		SurgicalMaskProtect:      0.85,
		FFPXMaskProtect:          0.99,
		HydroAlcoholicGelProtect: 0.99,
		Parameters: []*covidtracker.RiskParameter{
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.Aircraft, Duration: covidtracker.Long},
				NbDirect:               10,
				ProbaContagionDirect:   0.7,
				NbContact:              3,
				ProbaContagionContact:  0.6,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{"Segment avion long"},
				Advices:                []string{},
				MaskProtectDirect:      0.75,
				MaskProtectContact:     0.10,
				GelProtectContact:      0.90,
				MaskProtectIndirect:    0.10,
				GelProtectIndirect:     0.90,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.Aircraft, Duration: covidtracker.Normal},
				NbDirect:               7,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.6,
				NbIndirect:             200,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{},
				Advices:                []string{},
				MaskProtectDirect:      0.75,
				MaskProtectContact:     0.10,
				GelProtectContact:      0.90,
				MaskProtectIndirect:    0.10,
				GelProtectIndirect:     0.90,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.Aircraft, Duration: covidtracker.Short},
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.6,
				NbIndirect:             75,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{"Segment avion court"},
				Minuses:                []string{},
				Advices:                []string{},
				MaskProtectDirect:      0.75,
				MaskProtectContact:     0.10,
				GelProtectContact:      0.90,
				MaskProtectIndirect:    0.10,
				GelProtectIndirect:     0.90,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.TGV, Duration: covidtracker.Long},
				NbDirect:               5,
				ProbaContagionDirect:   0.8,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             50,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{"Segment train long"},
				Advices:                []string{},
				MaskProtectDirect:      0.75,
				MaskProtectContact:     0.10,
				GelProtectContact:      0.90,
				MaskProtectIndirect:    0.10,
				GelProtectIndirect:     0.90,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.TGV, Duration: covidtracker.Normal},
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.4,
				NbIndirect:             40,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{},
				Advices:                []string{},
				MaskProtectDirect:      0.85,
				MaskProtectContact:     0.15,
				GelProtectContact:      0.90,
				MaskProtectIndirect:    0.15,
				GelProtectIndirect:     0.90,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.TGV, Duration: covidtracker.Short},
				NbDirect:               5,
				ProbaContagionDirect:   0.6,
				NbContact:              2,
				ProbaContagionContact:  0.4,
				NbIndirect:             40,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{"Segment train court"},
				Minuses:                []string{},
				Advices:                []string{},
				MaskProtectDirect:      0.90,
				MaskProtectContact:     0.15,
				GelProtectContact:      0.90,
				MaskProtectIndirect:    0.15,
				GelProtectIndirect:     0.90,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.TER, Duration: covidtracker.Long},
				NbDirect:               5,
				ProbaContagionDirect:   0.8,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             50,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{"Segment train long"},
				Advices:                []string{},
				MaskProtectDirect:      0.75,
				MaskProtectContact:     0.10,
				GelProtectContact:      0.90,
				MaskProtectIndirect:    0.10,
				GelProtectIndirect:     0.90,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.TER, Duration: covidtracker.Normal},
				NbDirect:               5,
				ProbaContagionDirect:   0.7,
				NbContact:              2,
				ProbaContagionContact:  0.4,
				NbIndirect:             40,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{},
				Advices:                []string{},
				MaskProtectDirect:      0.85,
				MaskProtectContact:     0.15,
				GelProtectContact:      0.90,
				MaskProtectIndirect:    0.15,
				GelProtectIndirect:     0.90,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.TER, Duration: covidtracker.Short},
				NbDirect:               5,
				ProbaContagionDirect:   0.6,
				NbContact:              2,
				ProbaContagionContact:  0.4,
				NbIndirect:             40,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{"Segment train court"},
				Minuses:                []string{},
				Advices:                []string{},
				MaskProtectDirect:      0.90,
				MaskProtectContact:     0.15,
				GelProtectContact:      0.90,
				MaskProtectIndirect:    0.15,
				GelProtectIndirect:     0.90,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.CarSolo, Duration: covidtracker.Long},
				NbDirect:               2,
				ProbaContagionDirect:   0.4,
				NbContact:              0,
				ProbaContagionContact:  0.,
				NbIndirect:             8,
				ProbaContagionIndirect: 0.5,
				Pluses:                 []string{"Vous êtes seul(e) dans la voiture"},
				Minuses:                []string{"Segment voiture long, il faudra probablement s'arrêter à une pompe à essence"},
				Advices:                []string{"Lavez vous bien les mains si vous prenez de l'essence"},
				MaskProtectDirect:      0.,
				MaskProtectContact:     0.2,
				GelProtectContact:      0.9,
				MaskProtectIndirect:    0.2,
				GelProtectIndirect:     0.9,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.CarSolo, Duration: covidtracker.Normal},
				NbDirect:               0,
				ProbaContagionDirect:   0.,
				NbContact:              0,
				ProbaContagionContact:  0.,
				NbIndirect:             2,
				ProbaContagionIndirect: 0.5,
				Pluses:                 []string{"Vous êtes seul(e) dans la voiture"},
				Minuses:                []string{},
				Advices:                []string{"Lavez vous bien les mains si vous prenez de l'essence"},
				MaskProtectDirect:      0.,
				MaskProtectContact:     0.2,
				GelProtectContact:      0.9,
				MaskProtectIndirect:    0.2,
				GelProtectIndirect:     0.9,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.CarSolo, Duration: covidtracker.Short},
				NbDirect:               0,
				ProbaContagionDirect:   0.,
				NbContact:              0,
				ProbaContagionContact:  0.,
				NbIndirect:             2,
				ProbaContagionIndirect: 0.5,
				Pluses:                 []string{"Segment voiture court", "Vous êtes seul(e) dans la voiture"},
				Minuses:                []string{},
				Advices:                []string{},
				MaskProtectDirect:      0.,
				MaskProtectContact:     0.2,
				GelProtectContact:      0.9,
				MaskProtectIndirect:    0.2,
				GelProtectIndirect:     0.9,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.CarDuo, Duration: covidtracker.Long},
				NbDirect:               1,
				ProbaContagionDirect:   0.9,
				NbContact:              1,
				ProbaContagionContact:  0.6,
				NbIndirect:             8,
				ProbaContagionIndirect: 0.5,
				Pluses:                 []string{},
				Minuses:                []string{"Segment voiture long, il faudra probablement s'arrêter à une pompe à essence", "Vous êtes plusieurs dans voiture"},
				Advices:                []string{"Lavez vous bien les mains si vous prenez de l'essence"},
				MaskProtectDirect:      0.7,
				MaskProtectContact:     0.2,
				GelProtectContact:      0.6,
				MaskProtectIndirect:    0.2,
				GelProtectIndirect:     0.6,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.CarDuo, Duration: covidtracker.Normal},
				NbDirect:               1,
				ProbaContagionDirect:   0.85,
				NbContact:              1,
				ProbaContagionContact:  0.6,
				NbIndirect:             5,
				ProbaContagionIndirect: 0.3,
				Pluses:                 []string{},
				Minuses:                []string{"Vous êtes plusieurs dans voiture"},
				Advices:                []string{"Lavez vous bien les mains si vous prenez de l'essence"},
				MaskProtectDirect:      0.75,
				MaskProtectContact:     0.2,
				GelProtectContact:      0.6,
				MaskProtectIndirect:    0.2,
				GelProtectIndirect:     0.6,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.CarDuo, Duration: covidtracker.Short},
				NbDirect:               1,
				ProbaContagionDirect:   0.8,
				NbContact:              1,
				ProbaContagionContact:  0.6,
				NbIndirect:             5,
				ProbaContagionIndirect: 0.3,
				Pluses:                 []string{"Segment voiture court"},
				Minuses:                []string{"Vous êtes plusieurs dans voiture"},
				Advices:                []string{},
				MaskProtectDirect:      0.8,
				MaskProtectContact:     0.2,
				GelProtectContact:      0.6,
				MaskProtectIndirect:    0.2,
				GelProtectIndirect:     0.6,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.CarGroup, Duration: covidtracker.Long},
				NbDirect:               3,
				ProbaContagionDirect:   0.9,
				NbContact:              4,
				ProbaContagionContact:  0.6,
				NbIndirect:             10,
				ProbaContagionIndirect: 0.5,
				Pluses:                 []string{},
				Minuses:                []string{"Segment voiture long, il faudra probablement s'arrêter à une pompe à essence", "Vous êtes plusieurs dans voiture"},
				Advices:                []string{"Lavez vous bien les mains si vous prenez de l'essence"},
				MaskProtectDirect:      0.7,
				MaskProtectContact:     0.2,
				GelProtectContact:      0.6,
				MaskProtectIndirect:    0.2,
				GelProtectIndirect:     0.6,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.CarGroup, Duration: covidtracker.Normal},
				NbDirect:               3,
				ProbaContagionDirect:   0.8,
				NbContact:              4,
				ProbaContagionContact:  0.6,
				NbIndirect:             5,
				ProbaContagionIndirect: 0.3,
				Pluses:                 []string{},
				Minuses:                []string{"Vous êtes plusieurs dans voiture"},
				Advices:                []string{"Lavez vous bien les mains si vous prenez de l'essence"},
				MaskProtectDirect:      0.75,
				MaskProtectContact:     0.2,
				GelProtectContact:      0.6,
				MaskProtectIndirect:    0.2,
				GelProtectIndirect:     0.6,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.CarGroup, Duration: covidtracker.Short},
				NbDirect:               3,
				ProbaContagionDirect:   0.75,
				NbContact:              3,
				ProbaContagionContact:  0.5,
				NbIndirect:             3,
				ProbaContagionIndirect: 0.2,
				Pluses:                 []string{"Segment voiture court"},
				Minuses:                []string{"Vous êtes plusieurs dans voiture"},
				Advices:                []string{},
				MaskProtectDirect:      0.8,
				MaskProtectContact:     0.2,
				GelProtectContact:      0.7,
				MaskProtectIndirect:    0.2,
				GelProtectIndirect:     0.7,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.PublicTransports, Duration: covidtracker.Long},
				NbDirect:               100,
				ProbaContagionDirect:   0.8,
				NbContact:              10,
				ProbaContagionContact:  0.8,
				NbIndirect:             300,
				ProbaContagionIndirect: 0.4,
				Pluses:                 []string{},
				Minuses:                []string{"Segment de transports en commun long"},
				Advices:                []string{},
				MaskProtectDirect:      0.6,
				MaskProtectContact:     0.1,
				GelProtectContact:      0.5,
				MaskProtectIndirect:    0.1,
				GelProtectIndirect:     0.5,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.PublicTransports, Duration: covidtracker.Normal},
				NbDirect:               50,
				ProbaContagionDirect:   0.8,
				NbContact:              5,
				ProbaContagionContact:  0.8,
				NbIndirect:             150,
				ProbaContagionIndirect: 0.4,
				Pluses:                 []string{},
				Minuses:                []string{},
				Advices:                []string{},
				MaskProtectDirect:      0.7,
				MaskProtectContact:     0.1,
				GelProtectContact:      0.5,
				MaskProtectIndirect:    0.1,
				GelProtectIndirect:     0.5,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.PublicTransports, Duration: covidtracker.Short},
				NbDirect:               10,
				ProbaContagionDirect:   0.7,
				NbContact:              3,
				ProbaContagionContact:  0.5,
				NbIndirect:             60,
				ProbaContagionIndirect: 0.4,
				Pluses:                 []string{},
				Minuses:                []string{},
				Advices:                []string{},
				MaskProtectDirect:      0.7,
				MaskProtectContact:     0.1,
				GelProtectContact:      0.6,
				MaskProtectIndirect:    0.1,
				GelProtectIndirect:     0.6,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.TaxiSolo, Duration: covidtracker.Normal},
				NbDirect:               1,
				ProbaContagionDirect:   0.85,
				NbContact:              1,
				ProbaContagionContact:  0.6,
				NbIndirect:             5,
				ProbaContagionIndirect: 0.3,
				Pluses:                 []string{},
				Minuses:                []string{"Vous êtes plusieurs dans voiture"},
				Advices:                []string{},
				MaskProtectDirect:      0.75,
				MaskProtectContact:     0.2,
				GelProtectContact:      0.6,
				MaskProtectIndirect:    0.2,
				GelProtectIndirect:     0.6,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.TaxiGroup, Duration: covidtracker.Normal},
				NbDirect:               3,
				ProbaContagionDirect:   0.8,
				NbContact:              4,
				ProbaContagionContact:  0.6,
				NbIndirect:             5,
				ProbaContagionIndirect: 0.3,
				Pluses:                 []string{},
				Minuses:                []string{"Vous êtes plusieurs dans voiture"},
				Advices:                []string{},
				MaskProtectDirect:      0.75,
				MaskProtectContact:     0.2,
				GelProtectContact:      0.6,
				MaskProtectIndirect:    0.2,
				GelProtectIndirect:     0.6,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.Scooter, Duration: covidtracker.Normal},
				NbDirect:               1,
				ProbaContagionDirect:   0.3,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             5,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{},
				Advices:                []string{},
				MaskProtectDirect:      0.7,
				MaskProtectContact:     0.1,
				GelProtectContact:      0.5,
				MaskProtectIndirect:    0.1,
				GelProtectIndirect:     0.5,
			},
			{
				Scope:                  covidtracker.ParameterScope{Transportation: covidtracker.Bike, Duration: covidtracker.Normal},
				NbDirect:               1,
				ProbaContagionDirect:   0.3,
				NbContact:              2,
				ProbaContagionContact:  0.5,
				NbIndirect:             5,
				ProbaContagionIndirect: 0.1,
				Pluses:                 []string{},
				Minuses:                []string{},
				Advices:                []string{},
				MaskProtectDirect:      0.7,
				MaskProtectContact:     0.1,
				GelProtectContact:      0.5,
				MaskProtectIndirect:    0.1,
				GelProtectIndirect:     0.5,
			},
		},
	}
	return dal.Insert(defaultParams)
}
