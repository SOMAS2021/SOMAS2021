package health

import (
	"math"

	"github.com/SOMAS2021/SOMAS2021/pkg/config"
)

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
	// HP loss when currentHP = weakLevel
	HPLossBase int
	// HP loss slope w.r.t currentHP - weakLevel
	HPLossSlope float64
}

func NewHealthInfo(parameters *config.ConfigParameters) *HealthInfo {
	return &HealthInfo{
		MaxHP:          parameters.MaxHP,
		WeakLevel:      parameters.WeakLevel,
		Width:          parameters.Width,
		Tau:            parameters.Tau,
		HPReqCToW:      parameters.HpReqCToW,
		HPCritical:     parameters.HpCritical,
		MaxDayCritical: parameters.MaxDayCritical,
		HPLossBase:     parameters.HPLossBase,
		HPLossSlope:    parameters.HPLossSlope,
	}
}

func reachHP(currentHP float64, goalHP float64, healthInfo *HealthInfo) int { // this function is from team6 strategies just made public so everyone can use it
	denom := healthInfo.Width - goalHP + (1-healthInfo.HPLossSlope)*currentHP - float64(healthInfo.HPLossBase) + healthInfo.HPLossSlope*float64(healthInfo.WeakLevel)
	return int(healthInfo.Tau * math.Log(healthInfo.Width/denom))
}

func maintainHP(currentHP float64, healthInfo *HealthInfo) int { // this function is a modified verion of the function from team6 strategies just made public so everyone can use it
	goalHP := float64(currentHP)
	denom := healthInfo.Width - goalHP + (1-healthInfo.HPLossSlope)*currentHP - float64(healthInfo.HPLossBase) + healthInfo.HPLossSlope*float64(healthInfo.WeakLevel)
	return int(healthInfo.Tau * math.Log(healthInfo.Width/denom))
}
