package agents

import (
	"errors"
	"log"
<<<<<<< HEAD:pkg/agents/baseagent.go

	"github.com/SOMAS2021/SOMAS2021/pkg/infra/tower"
=======
	messages "github.com/SOMAS2021/SOMAS2021/pkg/infra/messages"
	tower "github.com/SOMAS2021/SOMAS2021/pkg/infra/tower"
>>>>>>> 4bbf3a2e425fa4818564f481cfc3d4cd85548358:pkg/agents/default/baseagent.go
	"github.com/divan/goabm/abm"
	"github.com/google/uuid"
)

type Agent interface {
	Run()
}

type Base struct {
	hp    int
	floor int
	id    string
	tower *tower.Tower
	inbox chan messages.Message
}

func NewBaseAgent(abm *abm.ABM, floor, hp int) (*Base, error) {
	world := abm.World()
	if world == nil {
		return nil, errors.New("Agent needs a World defined to operate")
	}
	tower, ok := world.(*tower.Tower)
	if !ok {
		return nil, errors.New("Agent needs a Tower world to operate")
	}
	return &Base{
		floor: floor,
		hp:    hp,
		tower: tower,
		id:    uuid.New().String(),
	}, nil
}

func (a *Base) Run() {
	log.Printf("An agent cycle executed from base agent %d", a.floor)
}

func (a *Base) HP() int {
	return a.hp
}

func (a *Base) Floor() int {
	return a.floor
}

<<<<<<< HEAD:pkg/agents/baseagent.go
func (a *Base) ID() string {
	return a.id
=======
func (a *BaseAgent) recieveMessage() msg messages.Message{
    msg := <- a.inbox
    return msg
>>>>>>> 4bbf3a2e425fa4818564f481cfc3d4cd85548358:pkg/agents/default/baseagent.go
}
