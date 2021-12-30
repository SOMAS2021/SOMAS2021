package team2

func calcReward(hp int, hpInc int) float64 {
	//TODO: reward should be redone later according to new healthInfo
	ret := float64(0.0)

	//we encourage agent to survive
	if hp > 20 {
		ret += 1.0
	} else {
		ret -= 0.5
	}
	// Use slop instead of hard threshold?
	threshold := float64(20)
	hpT := float64(hp)
	if hpT > threshold {
		ret += 1.0*(1.0/(100.0-threshold))*hpT - threshold/(100.0-threshold)
	} else {
		ret -= 0.5 * ((-1.0/(threshold))*hpT + 1)
	}

	//we encourage ageny to eat less when hp level is high
	oldHP := float64(hp - hpInc)
	incRate := float64(hpInc) / oldHP
	if incRate > 1.0 {
		return ret + 1.0
	}
	return ret - 1.0
	//TODO: we encourage agent to save other agent
	// if actionTaken == 0 => ret+= 1 * num_of_saved_agent
}

func (a *CustomAgent2) updateRTable(hpInc int, state int, action int) {
	reward := calcReward(a.HP(), hpInc)
	a.rTable[state][action] = reward
}
