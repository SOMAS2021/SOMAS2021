package day

type DayInfo struct {
	TicksPerFloor          int
	TicksPerDay            int
	SimulationDays         int
	DaysPerReshuffle       int
	TicksPerReshuffle      int
	TotalTicks             int
	CurrTick               int
	CurrDay                int
	BehaviourCtr           map[string]int
	BehaviourCtrData       [][]string
	BehaviourChangeCtr     map[string]int
	BehaviourChangeCtrData [][]string
	Utility                float64 // Convert to string before appending to UtilityData
	UtilityData            [][]string
	DeathData              [][]string
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
		BehaviourCtr: map[string]int{
			"Altruist":     0,
			"Collectivist": 0,
			"Selfish":      0,
			"Narcissist":   0,
		},
		BehaviourCtrData: [][]string{},
		BehaviourChangeCtr: map[string]int{
			"A2A": 0,
			"A2C": 0,
			"A2S": 0,
			"A2N": 0,
			"C2A": 0,
			"C2C": 0,
			"C2S": 0,
			"C2N": 0,
			"S2A": 0,
			"S2C": 0,
			"S2S": 0,
			"S2N": 0,
			"N2A": 0,
			"N2C": 0,
			"N2S": 0,
			"N2N": 0,
		},
		BehaviourChangeCtrData: [][]string{},
		Utility:                0.0,
		UtilityData:            [][]string{},
		DeathData:              [][]string{},
	}
}
