package agentTrust

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
)

type CustomAgent4 struct {
	*infra.Base
	myNumber int
	globalTrust float32
	globalTrustAdd float32
	globalTrustSubtract float32
	coefficients []float32
	lastFoodTaken int
	IntendedFoodTaken int 
}

func New(baseAgent *infra.Base) (agent.Agent, error) {
	return &CustomAgent4{
		Base:     baseAgent,
		
		globalTrust: 0.0, // MAKE SURE TO AMEND FIGURES FOR SENSIBLE AGENT BEHAVIOUR
		globalTrustAdd: 9.0,
		globalTrustSubtract: -9.0,
		coefficients: [0.1],
		
		lastFoodTaken: 0,
		IntendedFoodTaken: 0
	}, nil
}

func (a *CustomAgent4) Run() {
	a.Log("Reporting agent state", infra.Fields{"health": a.HP(), "floor": a.Floor()})

		// receivedMsg := a.Base.ReceiveMessage()
		// switch receivedMsg.MessageType() {
		// case "AckMessage":
		// 	a.globalTrust += a.globalTrustAdd  *  coefficients[0]// TODO AND WORK ON
		// 	if a.globalTrust > 100.0{
		// 		a.globalTrust = 100.0
		// 	}
		// // case "foodOnPlatMessage":
		// // 	if receivedMsg.food == a.CurrPlatFood() && a.CurrPlatFood() != -1
	 	//  case "LeaveFoodMessage":
		// 	if receivedMsg.food == a.currPlatFood() && receivedMsg.senderFloor - a.Floor() == -1 && a.CurrPlatFood() != -1{ // on the floor above you
		// 		a.globalTrust+= a.globalTrustAdd //
		// 		if a.globalTrust > 100.0{
		// 			a.globalTrust = 100.0
		// 		}

		// 	} else if receivedMsg.food != a.currPlatFood() && receivedMsg.senderFloor -a.Floor() == 1 && a.CurrPlatFood() != -1{ // on the floor below you

		// 	}


		// default:

		// }
	receivedMsg := a.ReceiveMessage()
	if receivedMsg != nil {
		receivedMsg.Visit(a)
	} else {
		a.Log("I got nothing")
	}

	a.IntendedFoodTaken = food.FoodType(a.currPlatFood() * (1- a.globalTrust/100))
	a.lastFoodTaken = a.TakeFood(a.IntendedFoodTaken)
	
}

func (a *CustomAgent2) HandleAskHP(msg messages.AskHPMessage) {
	reply := msg.Reply(a.Floor(), a.HP())
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I received an askHP message from ", infra.Fields{"floor": msg.SenderFloor(), "hp": reply.hp}})
}

func (a *CustomAgent2) HandleAskFoodTaken(msg messages.AskFoodTakenMessage) {
	reply := msg.Reply(a.Floor(), a.lastFoodTaken)
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I received an askFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": reply.food})
}

func (a *CustomAgent2) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {
	reply := msg.Reply(a.Floor(), a.IntendedFoodTaken)
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I received an askIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": reply.food}})
}

func (a *CustomAgent2) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	reply := msg.Reply(a.Floor(), true) // Change for later dependent on circumstance
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I received a requestLeaveFood message from ", infra.Fields{"floor": msg.SenderFloor()}})
}

func (a *CustomAgent2) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	reply := msg.Reply(a.Floor(), true) // Change for later dependent on circumstance
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I received a requestTakeFood message from ", infra.Fields{"floor": msg.SenderFloor()}})
}

func (a *CustomAgent2) HandleResponse(msg messages.BoolResponseMessage) {
	response := msg.Response() // Change for later dependent on circumstance
	a.Log("I received a Response message from ", infra.Fields{"floor": msg.SenderFloor(), "response": response})
}

func (a *CustomAgent2) HandleStateFoodTaken(msg messages.StateFoodTakenMessage) {
	statement := msg.Statement()
	a.Log("I received a StateFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": statement}})
}

func (a *CustomAgent2) HandleStateHP(msg messages.StateHPMessage) {
	statement := msg.Statement()
	a.Log("I received a StateHP message from ", infra.Fields{"floor": msg.SenderFloor(), "hp": statement})
}

func (a *CustomAgent2) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {
	statement := msg.Statement()
	a.Log("I received a StateIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": statement})
}