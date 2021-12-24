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
	//we encourage ageny to eat less
	if hpInc <= 20 {
		ret += 1.0
	} else {
		ret -= 0.5
	}
	return ret
	//TODO: we encourage agent to save other agent
	// if actionTaken == 0 => ret+= 1 * num_of_saved_agent
}

func (a *CustomAgent2) updateRTable(hpInc int, state int, action int) {
	reward := calcReward(a.HP(), hpInc)
	a.rTable[state][action] = reward
}
