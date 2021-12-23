package health

type HealthInfo struct {
	// Maximum HP
	MaxHP int
	// Base HP for each health level
	StrongLevel  int
	HealthyLevel int
	WeakLevel    int
	// Parameters of the updateHP function
	Width float64
	Tau   float64
	// Parameters of the hpDecay function - Costs of living
	CostStrong  int
	CostHealthy int
	CostWeak    int
	// HP required to leave the critical level and reach the weak level
	HPReqCToW int
	// HP value attributed in the critical level
	HPCritical int
	// Number of days an agent can stay critical before dying
	MaxDayCritical int
}

func NewHealthInfo(MaxHP, StrongLevel, HealthyLevel, WeakLevel int, Width, Tau float64, CostStrong, CostHealthy, CostWeak, HPReqCToW, HPCritical, MaxDayCritical int) *HealthInfo {
	return &HealthInfo{
		MaxHP:          MaxHP,
		StrongLevel:    StrongLevel,
		HealthyLevel:   HealthyLevel,
		WeakLevel:      WeakLevel,
		Width:          Width,
		Tau:            Tau,
		CostStrong:     CostStrong,
		CostHealthy:    CostHealthy,
		CostWeak:       CostWeak,
		HPReqCToW:      HPReqCToW,
		HPCritical:     HPCritical,
		MaxDayCritical: MaxDayCritical,
	}
}
