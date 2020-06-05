package covidtracker

type Transportation string

const (
	TGV              Transportation = "tgv"
	TER              Transportation = "ter"
	Aircraft         Transportation = "aircraft"
	Car              Transportation = "car"
	CarSolo          Transportation = "car-solo"
	CarDuo           Transportation = "car-duo"
	CarGroup         Transportation = "car-group"
	TaxiSolo         Transportation = "taxi-solo"
	TaxiGroup        Transportation = "taxi-group"
	PublicTransports Transportation = "public-transports"
	Scooter          Transportation = "scooter"
	Bike             Transportation = "bike"
)
