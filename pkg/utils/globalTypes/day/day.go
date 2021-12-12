package day

type DayInfo struct {
	TicksPerDay       int
	SimulationDays    int
	DaysPerReshuffle  int
	TicksPerReshuffle int
	TotalTicks        int
	CurrTick          int
}

func NewDayInfo(TicksPerDay, SimulationDays, DaysPerReshuffle int) *DayInfo {
	return &DayInfo{
		TicksPerDay:       TicksPerDay,
		SimulationDays:    SimulationDays,
		DaysPerReshuffle:  DaysPerReshuffle,
		TicksPerReshuffle: TicksPerDay * DaysPerReshuffle,
		TotalTicks:        TicksPerDay * SimulationDays,
		CurrTick:          1,
	}
}
