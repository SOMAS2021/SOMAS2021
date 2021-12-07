package tower

type BaseAgentCore struct {
	hp    int
	floor int
}

func (tower *Tower) GetHP(id string) int {
	return tower.agents[id].hp
}

func (tower *Tower) GetFloor(id string) int {
	return tower.agents[id].floor
}

func (tower *Tower) Exists(id string) bool {
	if _, found := tower.agents[id]; found {
		return true
	} else {
		return false
	}
}

func (tower *Tower) setFloor(id string, newFloor int) {
	temp := BaseAgentCore{
		hp:    tower.agents[id].hp,
		floor: newFloor,
	}
	tower.agents[id] = temp
}

func (tower *Tower) setHP(id string, newHP int) {
	temp := BaseAgentCore{
		hp:    newHP,
		floor: tower.agents[id].floor,
	}
	tower.agents[id] = temp
}

func (t *Tower) SetAgent(agentHP, agentFloor int, id string) {
	t.mx.Lock()
	t.agents[id] = BaseAgentCore{ // creating a new instance of agent in hash map
		hp:    agentHP,
		floor: agentFloor,
	}
	t.mx.Unlock()
}
