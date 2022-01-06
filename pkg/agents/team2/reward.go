package team2

import (
	"math"
)

func calcReward(hp int, hpInc int, MaxHp int, WeakLevel int, HPCritical int, DaysAtCritical int, MaxDayCritical int) float64 {
	//TODO: reward should be redone later according to new healthInfo
	ret := 0.0
	overEating := 86
	scale := 1.0

	//we encourage agent to survive
	/*if hp > 20 {
		ret += 1.0
	} else {
		ret -= 0.5
	}*/
	// Use slop instead of hard threshold?
	/*
		threshold := float64(20)
		hpT := float64(hp)
		if hpT > threshold {
			ret += 1.0*(1.0/(100.0-threshold))*hpT - threshold/(100.0-threshold)
		} else {
			ret -= 0.5 * ((-1.0/(threshold))*hpT + 1)
		}
	*/
	// Reward Cal base on HP + DaysAtCritical

	if hp > overEating { // over 2 tau: less reward for agent have more than enough food
		ret += scale * (math.Exp(-float64(hp-86) / 100))
	} else if hp >= WeakLevel { // more reward for agent in normal state
		ret += scale * (((1.0-0.3)/(100.0-float64(WeakLevel)))*float64(hp) - (1.0-0.3)*float64(WeakLevel)/(100.0-float64(WeakLevel)))
	} else { // penaty for agent at crital HP
		ret += -scale
	}
	if DaysAtCritical > 0 { // penaty for days at critical
		ret += 0.5 * scale * ((-1.0 / float64(MaxDayCritical)) * float64(DaysAtCritical))
	}
	if hp <= 0 { // penalty for death
		ret -= 3 * scale
	}

	return ret
	/*we encourage ageny to eat less when hp level is high
	oldHP := float64(hp - hpInc)
	incRate := float64(hpInc) / oldHP
	if incRate > 1.0 {
		return ret + 1.0
	}
	return ret - 1.0
	//TODO: we encourage agent to save other agent
	// if actionTaken == 0 => ret+= 1 * num_of_saved_agent
	*/
}

func (a *CustomAgent2) updateRTable(hpInc int, state int, action int) {
	reward := calcReward(a.HP(), hpInc, a.HealthInfo().MaxHP, a.HealthInfo().WeakLevel, a.HealthInfo().HPCritical, a.DaysAtCritical(), a.HealthInfo().MaxDayCritical)
	a.rTable[state][action] = reward
}
