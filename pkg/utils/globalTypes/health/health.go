package health

type HealthInfo struct {
	StrongLevel    int // defined as the lowest HP value for a given HP level
	HealthyLevel   int
	WeakLevel      int
	CriticalLevel  int
	FoodReqStrong  int // food requirement to stay at a given level
	FoodReqHealthy int
	FoodReqWeak    int
	FoodReqHToS    int // food requirement to switch from healthy to strong
	FoodReqWToH    int // food requirement to switch from weak to healthy
	FoodReqCToW    int // food requirement to switch from critical to weak
	MaxDayCritical int // maximum number of day at critical level
}

func NewHealthInfo(StrongLevel, HealthyLevel, WeakLevel, CriticalLevel, FoodReqStrong, FoodReqHealthy, FoodReqWeak, MaxDayCritical, FoodReqHToS, FoodReqWToH, FoodReqCToW int) *HealthInfo {
	return &HealthInfo{
		StrongLevel:    StrongLevel,
		HealthyLevel:   HealthyLevel,
		WeakLevel:      WeakLevel,
		CriticalLevel:  CriticalLevel,
		FoodReqStrong:  FoodReqStrong,
		FoodReqHealthy: FoodReqHealthy,
		FoodReqWeak:    FoodReqWeak,
		MaxDayCritical: MaxDayCritical,
		FoodReqHToS:    FoodReqHToS,
		FoodReqWToH:    FoodReqWToH,
		FoodReqCToW:    FoodReqCToW,
	}
}
