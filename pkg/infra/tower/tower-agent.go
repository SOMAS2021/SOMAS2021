package tower

type BaseAgentCore struct {
	hp        int
	floor     int
	agentType int
	inbox chan messages.Message
}

func (tower *Tower) GetHP(id string) int {
	tower.mx.RLock()
	defer tower.mx.RUnlock()
	return tower.agents[id].hp
}

func (tower *Tower) GetFloor(id string) int {
	tower.mx.RLock()
	defer tower.mx.RUnlock()
	return tower.agents[id].floor
}

func (tower *Tower) Exists(id string) bool {
	tower.mx.RLock()
	defer tower.mx.RUnlock()
	if _, found := tower.agents[id]; found {
		return true
	} else {
		return false
	}
}

func (tower *Tower) setFloor(id string, hp, newFloor, aType int) {
	tower.mx.Lock()
	tower.agents[id] = BaseAgentCore{
		hp:        hp,
		floor:     newFloor,
		agentType: aType,
	}
	tower.mx.Unlock()
}

func (tower *Tower) setHP(id string, newHP, floor, aType int) {
	tower.mx.Lock()
	tower.agents[id] = BaseAgentCore{
		hp:        newHP,
		floor:     1,
		agentType: 1,
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
