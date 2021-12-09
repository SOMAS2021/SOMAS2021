package tower

import(
	// "github.com/SOMAS2021/SOMAS2021/pkg/infra/messages"
	"container/list"
)

type BaseAgentCore struct {
	hp        int
	floor     int
	agentType int
	inbox 	  *list.List //chan messages.Message
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

func (tower *Tower) setFloor(id string, newFloor int) {
	tower.mx.Lock()
	agent := tower.agents[id]
	agent.floor = newFloor
	tower.agents[id] = agent
	tower.mx.Unlock()
}

func (tower *Tower) setHP(id string, newHP int) {
	tower.mx.Lock()
	agent := tower.agents[id]
	agent.hp = newHP
	tower.agents[id] = agent
	tower.mx.Unlock()
}

func (t *Tower) SetAgent(aType, agentHP, agentFloor int, id string) {
	t.mx.Lock()
	t.agents[id] = BaseAgentCore{
		hp:        agentHP,
		floor:     agentFloor,
		agentType: aType,
		inbox: list.New(),
	}
	t.mx.Unlock()
}
