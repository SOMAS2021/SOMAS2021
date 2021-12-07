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
	// TODO: fix concurrency issues
	hp := tower.agents[id].hp
	aType := tower.agents[id].agentType
	tower.mx.Lock()
	tower.agents[id] = BaseAgentCore{
		hp:        hp,
		floor:     newFloor,
		agentType: aType,
	}
	tower.mx.Unlock()
}

func (tower *Tower) setHP(id string, newHP int) {
	// TODO: fix concurrency issues
	floor := tower.agents[id].floor
	aType := tower.agents[id].agentType
	tower.mx.Lock()
	tower.agents[id] = BaseAgentCore{
		hp:        newHP,
		floor:     floor,
		agentType: aType,
	}
	tower.mx.Unlock()
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
