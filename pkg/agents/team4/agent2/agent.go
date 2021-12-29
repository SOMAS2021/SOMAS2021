package agentTrust

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	"math/rand"
)

type CustomAgent4 struct {
	*infra.Base
	globalTrust float32
	globalTrustAdd float32
	globalTrustSubtract float32
	coefficients []float32
	lastFoodTaken food.FoodType
	IntendedFoodTaken food.FoodType
	sentMessages MessageMemory
	MessageToSend int
}

type MessageMemory struct{
	// An array of messages stored into the agent before platform change.
	//direction []int
	direction []int
	messages []messages.Message
}
func (a *CustomAgent4) sendingMessage(){
	switch a.MessageToSend % 15  {
		case 0:
			msg := messages.NewAskFoodTakenMessage(a.Floor())
			a.SendMessage(1, msg)
			a.sentMessages.direction = append(a.sentMessages.direction, 1)
			a.sentMessages.messages = append(a.sentMessages.messages, msg)
			a.Log("I sent a message", infra.Fields{"message": "AskFoodTaken"})
		case 1:
			msg := messages.NewAskFoodTakenMessage(a.Floor())
			a.SendMessage(-1, msg)
			a.sentMessages.direction = append(a.sentMessages.direction, -1)
			a.sentMessages.messages = append(a.sentMessages.messages, msg)
			a.Log("I sent a message", infra.Fields{"message": "AskFoodTaken"})
		case 2:
			msg := messages.NewAskHPMessage(a.Floor())
			a.SendMessage(1, msg)
			a.sentMessages.direction = append(a.sentMessages.direction, 1)
			a.sentMessages.messages = append(a.sentMessages.messages, msg)
			a.Log("I sent a message", infra.Fields{"message": "AskHP"})
		case 3:
			msg := messages.NewAskHPMessage(a.Floor())
			a.SendMessage(-1, msg)
			a.sentMessages.direction = append(a.sentMessages.direction, -1)
			a.sentMessages.messages = append(a.sentMessages.messages, msg)
			a.Log("I sent a message", infra.Fields{"message": "AskHP"})
		case 4:
			msg := messages.NewAskIntendedFoodIntakeMessage(a.Floor())
			a.SendMessage(1, msg)
			a.sentMessages.direction = append(a.sentMessages.direction, 1)
			a.sentMessages.messages = append(a.sentMessages.messages, msg)
			a.Log("I sent a message", infra.Fields{"message": "AskIntendedFoodIntake"})
		case 5:
			msg := messages.NewAskIntendedFoodIntakeMessage(a.Floor())
			a.SendMessage(-1, msg)
			a.sentMessages.direction = append(a.sentMessages.direction, -1)
			a.sentMessages.messages = append(a.sentMessages.messages, msg)
			a.Log("I sent a message", infra.Fields{"message": "AskIntendedFoodIntake"})
		case 6:
			msg := messages.NewRequestLeaveFoodMessage(a.Floor(), 10) //need to change how much to request to leave
			a.SendMessage(1, msg)
			a.sentMessages.direction = append(a.sentMessages.direction, 1)
			a.sentMessages.messages = append(a.sentMessages.messages, msg)
			a.Log("I sent a message", infra.Fields{"message": "RequestLeaveFood"})
		case 7:
			msg := messages.NewRequestLeaveFoodMessage(a.Floor(), 10) //need to change how much to request to leave
			a.SendMessage(-1, msg)
			a.sentMessages.direction = append(a.sentMessages.direction, -1)
			a.sentMessages.messages = append(a.sentMessages.messages, msg)
			a.Log("I sent a message", infra.Fields{"message": "RequestLeaveFood"})
		case 8:
			msg := messages.NewRequestTakeFoodMessage(a.Floor(), 20) //need to change how much to request to take
			a.SendMessage(1, msg)
			a.sentMessages.direction = append(a.sentMessages.direction, 1)
			a.sentMessages.messages = append(a.sentMessages.messages, msg)
			a.Log("I sent a message", infra.Fields{"message": "RequestTakeFood"})
		case 9:
			msg := messages.NewRequestTakeFoodMessage(a.Floor(), 20) //need to change how much to request to take
			a.SendMessage(-1, msg)
			a.sentMessages.direction = append(a.sentMessages.direction, -1)
			a.sentMessages.messages = append(a.sentMessages.messages, msg)
			a.Log("I sent a message", infra.Fields{"message": "RequestTakeFood"})
		case 10:
			msg := messages.NewStateFoodTakenMessage(a.Floor(), int(a.lastFoodTaken), "")
			a.SendMessage(1, msg)
			a.sentMessages.direction = append(a.sentMessages.direction, 1)
			a.sentMessages.messages = append(a.sentMessages.messages, msg)
			a.Log("I sent a message", infra.Fields{"message": "StateFoodTaken"})
		case 11:
			msg := messages.NewStateFoodTakenMessage(a.Floor(), int(a.lastFoodTaken), "")
			a.SendMessage(-1, msg)
			a.sentMessages.direction = append(a.sentMessages.direction, -1)
			a.sentMessages.messages = append(a.sentMessages.messages, msg)
			a.Log("I sent a message", infra.Fields{"message": "StateFoodTaken"})
		case 12:
			msg := messages.NewStateHPMessage(a.Floor(), a.HP(), "")
			a.SendMessage(1, msg)
			a.sentMessages.direction = append(a.sentMessages.direction, 1)
			a.sentMessages.messages = append(a.sentMessages.messages, msg)
			a.Log("I sent a message", infra.Fields{"message": "StateHP"})
		case 13:
			msg := messages.NewStateHPMessage(a.Floor(), a.HP(),  "")
			a.SendMessage(-1, msg)
			a.sentMessages.direction = append(a.sentMessages.direction, -1)
			a.sentMessages.messages = append(a.sentMessages.messages, msg)
			a.Log("I sent a message", infra.Fields{"message": "StateHP"})
		case 14:
			msg := messages.NewStateIntendedFoodIntakeMessage(a.Floor(), int(a.IntendedFoodTaken), "")
			a.SendMessage(1, msg)
			a.sentMessages.direction = append(a.sentMessages.direction, 1)
			a.sentMessages.messages = append(a.sentMessages.messages, msg)
			a.Log("I sent a message", infra.Fields{"message": "StateIntendedFoodIntake"})
		case 15:
			msg := messages.NewStateIntendedFoodIntakeMessage(a.Floor(), int(a.IntendedFoodTaken), "")
			a.SendMessage(-1, msg)
			a.sentMessages.direction = append(a.sentMessages.direction, -1)
			a.sentMessages.messages = append(a.sentMessages.messages, msg)
			a.Log("I sent a message", infra.Fields{"message": "StateIntendedFoodIntake"})
	}
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	return &CustomAgent4{
		Base:     baseAgent,
		
		globalTrust: 0.0, // TODO: Amend values for correct agent behaviour
		globalTrustAdd: 9.0, // TODO: Amend values for correct agent behaviour
		globalTrustSubtract: -9.0, // TODO: Amend values for correct agent behaviour
		coefficients: []float32{0.1}, // TODO: Amend values for correct agent behaviour
		
		// Initialise the amount of food our agent intends to eat.
		IntendedFoodTaken: 0,
		// Initialise the actual food taken on the previous run.
		lastFoodTaken: 0,

		// Initialise Agents individual message memory
		sentMessages: MessageMemory{
			direction: []int{},
			messages: []messages.Message{},
		}, 
		// Define what message to send during a run.
		MessageToSend: rand.Intn(15),
	}, nil
}

func (a *CustomAgent4) Run() {
	a.Log("Reporting agent state", infra.Fields{"health": a.HP(), "floor": a.Floor()})

	receivedMsg := a.ReceiveMessage()
	if receivedMsg != nil {
		receivedMsg.Visit(a)
	} else {
		a.Log("I got nothing")
	}
	//TODO: Define a threshold limit for other agents to respond to our sent message.
	a.sendingMessage()
	// msg := messages.NewRequestLeaveFoodMessage(a.Floor(), 10) //need to change how much to request to leave
	// a.SendMessage(-1, msg)
	// a.sentMessages.direction = append(a.sentMessages.direction, -1)
	// a.sentMessages.messages = append(a.sentMessages.messages, msg)
	// a.Log("I sent a message", infra.Fields{"message": "RequestLeaveFood"})
	// Verbose explanation of calculating the intented food taken
	// trust_prop := 100 - int(a.globalTrust)
	// intversion := int(a.CurrPlatFood())
	// food_prop := intversion * trust_prop / 100
	// a.IntendedFoodTaken = food.FoodType(food_prop)
	a.IntendedFoodTaken = food.FoodType(int(int(a.CurrPlatFood()) * (100 - int(a.globalTrust))/100))

	a.lastFoodTaken = a.TakeFood(a.IntendedFoodTaken)
	// MessageToSend +=1
	a.MessageToSend +=rand.Intn(15)
}




func (a *CustomAgent4) HandleAskHP(msg messages.AskHPMessage) {
	reply := msg.Reply(a.Floor(), a.HP(), msg.ID())
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I received an askHP message from ", infra.Fields{"floor": msg.SenderFloor(), "hp": a.HP()})
}

func (a *CustomAgent4) HandleAskFoodTaken(msg messages.AskFoodTakenMessage) {
	reply := msg.Reply(a.Floor(), int(a.lastFoodTaken), msg.ID())
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I received an askFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": a.lastFoodTaken})
}

func (a *CustomAgent4) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {
	reply := msg.Reply(a.Floor(), int(a.IntendedFoodTaken), msg.ID())
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I received an askIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": a.IntendedFoodTaken})
}

func (a *CustomAgent4) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	//fmt.Printf(msg.ID())
	reply := msg.Reply(a.Floor(), true, msg.ID()) // TODO: Change for later dependent on circumstance
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I received a requestLeaveFood message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgent4) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	reply := msg.Reply(a.Floor(), true, msg.ID()) // TODO: Change for later dependent on circumstance
	a.SendMessage(msg.SenderFloor()-a.Floor(), reply)
	a.Log("I received a requestTakeFood message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgent4) HandleResponse(msg messages.BoolResponseMessage) {
	response := msg.Response() // TODO: Change for later dependent on circumstance
	if !response{
	 	a.globalTrust-= a.globalTrustSubtract * a.coefficients[0] // TODO: adapt for other conditions
	}else{// Iterating through all messages in agent memory
		for i:=0; i< len(a.sentMessages.messages); i++{ //TODO: adapt when uuids are implemented so you would have reponse.uuid == sentmessage.uuid
			if msg.ID() == a.sentMessages.messages[i].ID(){
				a.Log("Received a message ", infra.Fields{"sender_uuid": msg.ID(), "sentmessage_uuid": a.sentMessages.messages[i].ID()})
			}
		}
	}
	//a.Log("I received a Response message from ", infra.Fields{"floor": msg.SenderFloor(), "response": response})
}

func (a *CustomAgent4) HandleStateFoodTaken(msg messages.StateFoodTakenMessage) {
	statement := msg.Statement()

	a.Log("I received a StateFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": statement})
}

func (a *CustomAgent4) HandleStateHP(msg messages.StateHPMessage) {
	statement := msg.Statement()
	a.Log("I received a StateHP message from ", infra.Fields{"floor": msg.SenderFloor(), "hp": statement})
}

func (a *CustomAgent4) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {
	statement := msg.Statement()
	a.Log("I received a StateIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": statement})
}




// func (a *CustomAgent4) Run() {
// 	a.Log("Reporting agent state", infra.Fields{"health": a.HP(), "floor": a.Floor()})

// 		receivedMsg := a.Base.ReceiveMessage()
// 		switch receivedMsg.MessageType() {
// 		case "AckMessage":
// 			a.globalTrust += a.globalTrustAdd  *  coefficients[0]// TODO AND WORK ON
// 			if a.globalTrust > 100.0{
// 				a.globalTrust = 100.0
// 			}
// 		// case "foodOnPlatMessage":
// 		// 	if receivedMsg.food == a.CurrPlatFood() && a.CurrPlatFood() != -1
// 	 	 case "LeaveFoodMessage":
// 			if receivedMsg.food == a.currPlatFood() && receivedMsg.senderFloor - a.Floor() == -1 && a.CurrPlatFood() != -1{ // on the floor above you
// 				a.globalTrust+= a.globalTrustAdd //
// 				if a.globalTrust > 100.0{
// 					a.globalTrust = 100.0
// 				}

// 			} else if receivedMsg.food != a.currPlatFood() && receivedMsg.senderFloor -a.Floor() == 1 && a.CurrPlatFood() != -1{ // on the floor below you

// 			}


// 		default:

// 		}
// 	receivedMsg := a.ReceiveMessage()
// 	if receivedMsg != nil {
// 		receivedMsg.Visit(a)
// 	} else {
// 		a.Log("I got nothing")
// 	}

// 	// trust_prop := 100 - int(a.globalTrust)
// 	// intversion := int(a.CurrPlatFood())
// 	// food_prop := intversion * trust_prop / 100
// 	// a.IntendedFoodTaken = food.FoodType(food_prop)

// 	a.IntendedFoodTaken = food.FoodType(int(int(a.CurrPlatFood()) * (100 - int(a.globalTrust))/100))

// 	a.lastFoodTaken = a.TakeFood(a.IntendedFoodTaken)
	
// }