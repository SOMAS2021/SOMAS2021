package team5

import (
	"math"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
)

type Memory struct {
	trust             float64 // scale of -5 to 5, with -5 being least trustworthy and 5 being most trustworthy, 0 is neutral
	favour            float64 // e.g. generosity; scale of -5 to 5, with -5 being least favoured and 5 being most favoured, 0 is neutral
	daysSinceLastSeen int     // days since last interaction
}

type CustomAgent5 struct {
	*infra.Base
	selflishness      float64
	lastMeal          float64
	daysSinceLastMeal int
	currentAim        int
	satisfaction      int
	lastHP            int
	lastCriticalCount int
	daysAlive         int
	attemptToEat      bool
	memory            map[string]Memory
}

func New(baseAgent *infra.Base) (agent.Agent, error) {
	return &CustomAgent5{
		Base:              baseAgent,
		selflishness:      3.0,                 // of 0 to 3, with 3 being completely selflish, 0 being completely selfless
		lastMeal:          0,                   //Stores value of the last amount of food taken
		daysSinceLastMeal: 0,                   //Count of how many days since last eating
		currentAim:        0,                   //Scale of 0 to 2, 0 being willing to lose health, 1 being maintaining health, 2 being gaining health
		satisfaction:      0,                   //Scale of -3 to 3, with 3 being satisfied and unsatisfied
		lastHP:            0,                   //Remember last hp so we can see when a day has passed
		lastCriticalCount: 0,                   //Remember critical count so we can see when a day has passed
		daysAlive:         0,                   //Count how many days agent has been alive
		attemptToEat:      true,                //Variable needed to check if we have already attempted to eat on a day
		memory:            map[string]Memory{}, // Memory of other agents, key is agent id
	}, nil
}

func (a *CustomAgent5) NewMemory(id string) {
	a.memory[id] = Memory{
		trust:             0,
		favour:            0,
		daysSinceLastSeen: 0,
	}
}

func (a *CustomAgent5) UpdateAim() {
	switch {
	case a.selflishness >= 3:
		a.currentAim = 2 //If fully selfish always try to gain health
	case a.HP() > 80 && a.selflishness == 2:
		a.currentAim = 1 //Try to maintain health if near max health if mostly selfish
	case a.HP() > 80 && a.selflishness < 2:
		a.currentAim = 0 //Willing to lose health near max health if mostly or completely selfless
	case a.HP() > 50 && a.selflishness == 2:
		a.currentAim = 2 //Try to gain health if mostly selfish when above half health
	case a.HP() > 50 && a.selflishness == 1:
		a.currentAim = 1 //Try to maintain half health even if being mostly selfless
	case a.HP() > 50 && a.selflishness == 0:
		a.currentAim = 0 //Willing to lose health if being completely selfless
	case a.HP() > 10 && a.selflishness >= 1:
		a.currentAim = 2 //Try to gain health if less than half health and being anything but completely selfless
	default:
		a.currentAim = 1 //Default to maintain health if being completely selfless at less than half health
	}
}

func (a *CustomAgent5) UpdateSelfishness() {
	if a.satisfaction == 3 {
		a.selflishness--
	}
	if a.satisfaction < 0 {
		a.selflishness++
	}
	//The above is a basic implementation for now while messaging is not functional
	//Once messages are implemented this function will be dependent on our social network and treaties etc
}

func (a *CustomAgent5) FoodGain() float64 {
	return a.HealthInfo().Tau * 3
}

func (a *CustomAgent5) FoodMaintain() float64 {
	return a.HealthInfo().Tau * math.Log(1-(float64((a.HP()+30))/(3*a.HealthInfo().Width))) * -1
}

func (a *CustomAgent5) UpdateSatisfaction() {
	if a.HP() > 100 {
		a.satisfaction = 3
	}
	if a.lastMeal == 0 && a.satisfaction > -3 {
		a.satisfaction--
	}
	if a.HP() < 25 && a.satisfaction > -3 {
		a.satisfaction--
	}
	if a.HP() > 75 && a.satisfaction < 3 {
		a.satisfaction++
	}
}

func (a *CustomAgent5) GetMessages() {
	receivedMsg := a.Base.ReceiveMessage()
	for receivedMsg != nil {
		//some message processing depending on the type of the message

		receivedMsg = a.Base.ReceiveMessage() // receive next message in the inbox
	}
}

func (a *CustomAgent5) SendMessages() {
	//function that will send all messages we need to the other agents
}

func (a *CustomAgent5) DayPassed() {
	a.daysAlive++
	a.daysSinceLastMeal++
	//Also update daySinceLastSeen for memory here
}

func (a *CustomAgent5) Run() {
	a.Log("Reporting agent state of team 5 agent", infra.Fields{"health": a.HP(), "floor": a.Floor()})
	a.UpdateSelfishness()
	a.UpdateAim()
	attemptFood := 0.0
	if a.HP() < 10 {
		attemptFood = 1 //No point taking more food than 1 if in critical state as we will only reach 10hp with any amount of food > 0
	} else {
		if a.currentAim == 1 {
			attemptFood = a.FoodMaintain()
		} else if a.currentAim == 2 {
			attemptFood = a.FoodGain()
		}
	}
	if a.CurrPlatFood() != -1 && a.attemptToEat { //When platform reaches our floor and we haven't tried to eat, then try to eat
		a.lastMeal = a.TakeFood(attemptFood)
		if a.lastMeal > 0 {
			a.Log("Team 5 agent has taken food", infra.Fields{"amount": a.lastMeal})
			a.daysSinceLastMeal = 0
		}
		a.UpdateSatisfaction()
		a.attemptToEat = false
	}
	if (a.CurrPlatFood() == -1 && !a.attemptToEat) || a.CurrPlatFood() == 100 { //When platform has passed our floor, take that as a day passing and update parameters
		//Check for CurrPlatFood == 100 needed because of rare case when you are on bottom floor then reshuffled to top, the platform never passes you
		a.DayPassed()
		a.attemptToEat = true
	}
}
