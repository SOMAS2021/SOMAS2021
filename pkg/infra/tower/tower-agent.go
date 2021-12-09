package tower

import (
	// "github.com/SOMAS2021/SOMAS2021/pkg/infra/messages"
	"container/list"
)

type BaseAgentCore struct {
	hp        int
	floor     int
	agentType int
	inbox     *list.List //chan messages.Message
}

func (tower *Tower) HP(id string) int {
	tower.mx.Lock()
	defer tower.mx.Unlock()
	return tower.agents[id].hp
}

func (tower *Tower) Floor(id string) int {
	tower.mx.Lock()
	defer tower.mx.Unlock()
	return tower.agents[id].floor
}

func (tower *Tower) Exists(id string) bool {
	tower.mx.Lock()
	defer tower.mx.Unlock()
	_, found := tower.agents[id]
	return found
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

func (t *Tower) SetAgent(id string, agentHp int, agentFloor int, agentType int) {
	t.mx.Lock()
	defer t.mx.Unlock()
	t.agents[id] = BaseAgentCore{
		hp:        agentHp,
		floor:     agentFloor,
		agentType: agentType,
		inbox:     list.New(),
	}
}
