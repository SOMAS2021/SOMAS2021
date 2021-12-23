package health

type HealthInfo struct {
	// Maximum HP
	MaxHP int
	// Border between critical and weak state
	WeakLevel int
	// Parameters of the updateHP function
	Width float64
	Tau   float64
	// HP required to leave the critical level and reach the weak level
	HPReqCToW int
	// HP value attributed in the critical level
	HPCritical int
	// Number of days an agent can stay critical before dying
	MaxDayCritical int
}

func NewHealthInfo(MaxHP, WeakLevel int, Width, Tau float64, HPReqCToW, HPCritical, MaxDayCritical int) *HealthInfo {
	return &HealthInfo{
		MaxHP:          MaxHP,
		WeakLevel:      WeakLevel,
		Width:          Width,
		Tau:            Tau,
		HPReqCToW:      HPReqCToW,
		HPCritical:     HPCritical,
		MaxDayCritical: MaxDayCritical,
	}
}
