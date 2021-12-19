package health

type HealthInfo struct {
	// Base HP for each health level
	StrongLevel    int 
	HealthyLevel   int
	WeakLevel      int
	// Width of every health level
	WidthStrong    float64 
	WidthHealthy   float64
	WidthWeak      float64
	WidthCritical  float64
	// Food intake required to stay at a given level
	FoodReqStrong  float64 
	FoodReqHealthy float64
	FoodReqWeak    float64
	FoodReqCToW    float64
	// Number of days an agent can stay critical before dying
	MaxDayCritical int
}

func NewHealthInfo(StrongLevel, HealthyLevel, WeakLevel int, FoodReqStrong, FoodReqHealthy, FoodReqWeak, FoodReqCToW float64, MaxDayCritical int) *HealthInfo {
	return &HealthInfo{
		StrongLevel:    StrongLevel,
		HealthyLevel:   HealthyLevel,
		WeakLevel:      WeakLevel,
		WidthStrong:    float64(100 - StrongLevel),
		WidthHealthy:   float64(StrongLevel - HealthyLevel),
		WidthWeak:      float64(HealthyLevel - WeakLevel),
		WidthCritical:  float64(WeakLevel),
		FoodReqStrong:  FoodReqStrong, // is defined as the "time constant" of the 1st order system for updateFood
		FoodReqHealthy: FoodReqHealthy,
		FoodReqWeak:    FoodReqWeak,
		FoodReqCToW:    FoodReqCToW,
		MaxDayCritical: MaxDayCritical,
	}
}
