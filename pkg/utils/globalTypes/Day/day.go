package day_type

type DayInfo struct {
	IterationsPerDay       int
	SimulationDays         int
	DaysPerReshuffle       int
	IterationsPerReshuffle int
	TotalIterations        int
	CurrTick               int
}

func DayInformationNew() *DayInformation {
	return &DayInformation{
		IterationsPerDay:       0,
		SimulationDays:         0,
		DaysPerReshuffle:       0,
		IterationsPerReshuffle: 0,
		TotalIterations:        0,
		CurrTick:               0,
	}
}
