package day_type

type DayInfo struct {
	IterationsPerDay       int
	SimulationDays         int
	DaysPerReshuffle       int
	IterationsPerReshuffle int
	TotalIterations        int
	CurrTick               int
}

func NewDayInfo(IterationsPerDay, SimulationDays, DaysPerReshuffle int) *DayInfo {
	return &DayInformation{
		IterationsPerDay:       IterationsPerDay,
		SimulationDays:         SimulationDays,
		DaysPerReshuffle:       DaysPerReshuffle,
		IterationsPerReshuffle: IterationsPerDay * DaysPerReshuffle,
		TotalIterations:        IterationsPerDay * SimulationDays,
		CurrTick:               0,
	}
}
