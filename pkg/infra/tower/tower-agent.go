package tower

type BaseAgentCore struct {
	hp        int
	floor     int
	agentType int
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
		hp:        tower.agents[id].hp,
		floor:     newFloor,
		agentType: tower.agents[id].agentType,
	}
	tower.agents[id] = temp
}

func (tower *Tower) setHP(id string, newHP int) {
	temp := BaseAgentCore{
		hp:        newHP,
		floor:     tower.agents[id].floor,
		agentType: tower.agents[id].agentType,
	}
	tower.agents[id] = temp
}

func (t *Tower) SetAgent(aType, agentHP, agentFloor int, id string) {
	t.mx.Lock()
	t.agents[id] = BaseAgentCore{
		hp:        agentHP,
		floor:     agentFloor,
		agentType: aType,
	}
	t.mx.Unlock()
}
