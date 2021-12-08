package tower

type BaseAgentCore struct {
	hp        int
	floor     int
	agentType int
	inbox chan messages.Message
}

func (tower *Tower) HP(id string) int {
	tower.mx.RLock()
	defer tower.mx.RUnlock()
	return tower.agents[id].hp
}

func (tower *Tower) Floor(id string) int {
	tower.mx.RLock()
	defer tower.mx.RUnlock()
	return tower.agents[id].floor
}

func (tower *Tower) Exists(id string) bool {
	tower.mx.RLock()
	defer tower.mx.RUnlock()
	_, found := tower.agents[id]
	return found
}

func (tower *Tower) setFloor(id string, hp, newFloor, agentType int) {
	tower.mx.Lock()
	defer tower.mx.Unlock()
	tower.agents[id] = BaseAgentCore{
		hp:        hp,
		floor:     newFloor,
		agentType: agentType,
	}
}

func (tower *Tower) setHP(id string, newHP, floor, agentType int) {
	tower.mx.Lock()
	defer tower.mx.Unlock()
	tower.agents[id] = BaseAgentCore{
		hp:        newHP,
		floor:     1,
		agentType: 1,
	}
}

func (t *Tower) SetAgent(id string, agentHp int, agentFloor int, agentType int) {
	t.mx.Lock()
	defer t.mx.Unlock()
	t.agents[id] = BaseAgentCore{
		hp:        agentHp,
		floor:     agentFloor,
		agentType: agentType,
	}
}
