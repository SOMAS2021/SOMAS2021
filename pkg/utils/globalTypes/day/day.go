package day

type DayInfo struct {
	TicksPerFloor     int
	TicksPerDay       int
	SimulationDays    int
	DaysPerReshuffle  int
	TicksPerReshuffle int
	TotalTicks        int
	CurrTick          int
	CurrDay           int
}

func NewDayInfo(TicksPerFloor, TicksPerDay, SimulationDays, DaysPerReshuffle int) *DayInfo {
	return &DayInfo{
		TicksPerFloor:     TicksPerFloor,
		TicksPerDay:       TicksPerDay,
		SimulationDays:    SimulationDays,
		DaysPerReshuffle:  DaysPerReshuffle,
		TicksPerReshuffle: TicksPerDay * DaysPerReshuffle,
		TotalTicks:        TicksPerDay * SimulationDays,
		CurrTick:          1,
		CurrDay:           1,
	}
}
