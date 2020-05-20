package covidtracker

type AgeGroup string

const (
	AllAges AgeGroup = "0"
	Under15 AgeGroup = "A" // A	 moins de 15 ans
	Under44 AgeGroup = "B" // B	 15-44 ans
	Under64 AgeGroup = "C" // C	 45-64 ans
	Under74 AgeGroup = "D" // D	 65-74 ans
	Over75  AgeGroup = "E" // E	 75 et plus
)
