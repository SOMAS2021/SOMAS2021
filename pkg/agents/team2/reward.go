package team2

func calcReward(hp int, hpInc int) float32 {
	//TODO: reward should be redone later according to new healthInfo
	ret := float32(0.0)

	//we encourage agent to survive
	if hp > 20 {
		ret += 1.0
	} else {
		ret -= 0.5
	}
	//we encourage ageny to eat less when hp level is high
	oldHP := float32(hp - hpInc)
	incRate := float32(hpInc) / oldHP
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
