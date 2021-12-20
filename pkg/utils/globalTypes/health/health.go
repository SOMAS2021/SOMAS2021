package health

type HealthInfo struct {
	// Maximum HP
	MaxHP int
	// Base HP for each health level
	StrongLevel  int
	HealthyLevel int
	WeakLevel    int
	// Width of every health level
	WidthStrong   int
	WidthHealthy  int
	WidthWeak     int
	WidthCritical int
	// Food intake required to stay at a given level
	TauStrong   float64
	TauHealthy  float64
	TauWeak     float64
	FoodReqCToW float64
	// Number of days an agent can stay critical before dying
	MaxDayCritical int
}

func NewHealthInfo(MaxHP, StrongLevel, HealthyLevel, WeakLevel int, TauStrong, TauHealthy, TauWeak, FoodReqCToW float64, MaxDayCritical int) *HealthInfo {
	return &HealthInfo{
		MaxHP:          MaxHP,
		StrongLevel:    StrongLevel,
		HealthyLevel:   HealthyLevel,
		WeakLevel:      WeakLevel,
		WidthStrong:    MaxHP - StrongLevel,
		WidthHealthy:   StrongLevel - HealthyLevel,
		WidthWeak:      HealthyLevel - WeakLevel,
		WidthCritical:  WeakLevel,
		TauStrong:      TauStrong,
		TauHealthy:     TauHealthy,
		TauWeak:        TauWeak,
		FoodReqCToW:    FoodReqCToW,
		MaxDayCritical: MaxDayCritical,
	}
}
