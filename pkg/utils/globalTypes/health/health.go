package health

type HealthInfo struct {
	// Base HP for each health level
	StrongLevel    int 
	HealthyLevel   int
	WeakLevel      int
	WidthStrong    float64 // width of every level
	WidthHealthy   float64
	WidthWeak      float64
	WidthCritical  float64
	FoodReqStrong  float64 // food requirement to stay at a given level
	FoodReqHealthy float64
	FoodReqWeak    float64
	FoodReqCToW    float64
	MaxDayCritical int // maximum number of day at critical level
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
