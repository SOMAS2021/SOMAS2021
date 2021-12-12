package day

type DayInfo struct {
	TicksPerDay       int
	SimulationDays         int
	DaysPerReshuffle       int
	TicksPerReshuffle int
	TotalTicks        int
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
