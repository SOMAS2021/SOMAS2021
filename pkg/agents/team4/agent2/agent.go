package agentTrust

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
)

type CustomAgent4 struct {
	*infra.Base
	myNumber int
	globalTrust int
}

func New(baseAgent *infra.Base) (agent.Agent, error) {
	return &CustomAgent4{
		Base:     baseAgent,
		myNumber: 0,
		globalTrust: 0.0, // MAKE SURE TO AMEND FIGURES FOR SENSIBLE AGENT BEHAVIOUR
		globalTicks: 0,
		globalTrustAdd: 9,
		globalTrustSubtract: -9,
		coefficients: [],
	}, nil
}

func (a *CustomAgent4) Run() {
	a.Log("Reporting agent state", infra.Fields{"health": a.HP(), "floor": a.Floor()})

	defer a.globalTicks++ // count number of passed ticks/ days passed
	if globalTicks %2 ==0 && globalTicks !=0{
		receivedMsg := a.Base.ReceiveMessage()
		switch receivedMsg.MessageType() {
		case "AckMessage":
			a.globalTrust += a.globalTrustAdd  *  0.1// TODO AND WORK ON
		// case "foodOnPlatMessage":
		// 	if receivedMsg.food == a.CurrPlatFood() && a.CurrPlatFood() != -1
	  case "LeaveFoodMessage":
			if receivedMsg.food == a.currPlatFood() && receivedMsg.senderFloor - a.Floor() == -1 && a.CurrPlatFood() != -1{ // on the floor above you
				a.globalTrust+= a.globalTrustAdd //
				if a.globalTrust > 100.0{
					a.globalTrust = 100.0
				}

			} else if receivedMsg.food != a.currPlatFood() && receivedMsg.senderFloor -a.Floor() == 1 && a.CurrPlatFood() != -1{ // on the floor below you

			}


		default:

		}



	}

	foodtake_amount = int(a.currPlatFood() * (1- a.globalTrust/100))

	if foodtake_amount !=0{
		a.TakeFood(int(a.currPlatFood() * (1- a.globalTrust/100)))
	}

	// switch receivedMsg.MessageType() {
	// case condition:
	//
	// }
	if receivedMsg != nil {
		a.Log("I sent a message", infra.Fields{"message": receivedMsg.MessageType()})
	} else {
		a.Log("I got nothing")
	}

	if (a.myNumber)%2 == 0 {
		msg := *messages.NewAckMessage(int(a.Floor()), true)
		a.SendMessage(1, msg)
		a.Log("I sent a message", infra.Fields{"message": msg.MessageType()})
	} else {
		msg := *messages.NewBaseMessage(int(a.Floor()))
		a.SendMessage(1, msg)
	}
	a.Log("My agent is doing something", infra.Fields{"thing": "potatoe", "another_thing": "another potatoe"})
	a.TakeFood(10)
}
