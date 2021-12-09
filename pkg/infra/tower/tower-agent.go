package tower

type BaseAgentCore struct {
	hp        int
	floor     int
	agentType int
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
