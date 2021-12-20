package agentTrust

import (
	"math/rand"

	"fmt"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

type CustomAgent4 struct {
	*infra.Base
	globalTrust         float32
	globalTrustAdd      float32
	globalTrustSubtract float32
	coefficients        []float32
	lastFoodTaken       food.FoodType
	IntendedFoodTaken   food.FoodType
	sentMessages        MessageMemory
	MessageToSend       int
	lastPlatFood        food.FoodType
	maxFoodLimit        food.FoodType
}

type MessageMemory struct {
	// An array of messages stored into the agent before platform change.
	//direction []int
	direction []int
	messages  []messages.Message
}

func (a *CustomAgent4) AppendToMessageMemory(direction int, msg messages.Message) {
	a.sentMessages.direction = append(a.sentMessages.direction, direction)
	a.sentMessages.messages = append(a.sentMessages.messages, msg)
}

func (a *CustomAgent4) SendingMessage(direction int) {

	var msg messages.Message
	if direction == 0 {
		direction = 1
	} else {
		direction = -1
	}
	switch a.MessageToSend % 8 {
	case 0:
		msg = messages.NewAskFoodTakenMessage(a.ID(), a.Floor(), a.Floor()+direction)
	case 1:
		msg = messages.NewAskHPMessage(a.ID(), a.Floor(), a.Floor()+direction)
	case 2:
		msg = messages.NewAskIntendedFoodIntakeMessage(a.ID(), a.Floor(), a.Floor()+direction)
	case 3:
		msg = messages.NewRequestLeaveFoodMessage(a.ID(), a.Floor(), a.Floor()+direction, 10) //need to change how much to request to leave
	case 4:
		msg = messages.NewRequestTakeFoodMessage(a.ID(), a.Floor(), a.Floor()+direction, 20) //need to change how much to request to take
	case 5:
		msg = messages.NewStateFoodTakenMessage(a.ID(), a.Floor(), a.Floor()+direction, int(a.lastFoodTaken))
	case 6:
		msg = messages.NewStateHPMessage(a.ID(), a.Floor(), a.Floor()+direction, a.HP())
	case 7:
		msg = messages.NewStateIntendedFoodIntakeMessage(a.ID(), a.Floor(), a.Floor()+direction, int(a.IntendedFoodTaken))
	}

	a.SendMessage(msg)
	a.AppendToMessageMemory(direction, msg)
	a.Log("I sent a message", infra.Fields{"message": msg.MessageType()})

}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	return &CustomAgent4{
		Base: baseAgent,

		globalTrust:         0.0,                           // TODO: Amend values for correct agent behaviour
		globalTrustAdd:      9.0,                           // TODO: Amend values for correct agent behaviour
		globalTrustSubtract: -9.0,                          // TODO: Amend values for correct agent behaviour
		coefficients:        []float32{0.1, 0.2, 0.4, 0.5}, // TODO: Amend values for correct agent behaviour

		// Initialise the amount of food our agent intends to eat.
		IntendedFoodTaken: 0,
		// Initialise the actual food taken on the previous run.
		lastFoodTaken: 0,

		// Initialise Agents individual message memory
		sentMessages: MessageMemory{
			direction: []int{},
			messages:  []messages.Message{},
		},
		// Define what message to send during a run.
		MessageToSend: rand.Intn(8),
		lastPlatFood:  -1,
		maxFoodLimit:  50,
	}, nil
}

func (a *CustomAgent4) Run() {
	a.Log("Reporting agent state", infra.Fields{"health": a.HP(), "floor": a.Floor()})

	// if food.FoodType(a.CurrPlatFood()) != a.lastPlatFood && a.HasEaten(){ //TODO: Change if we don't eat everyday
	// 	a.lastPlatFood = a.CurrPlatFood()
	// }
	// if food.FoodType(a.CurrPlatFood()) != a.lastPlatFood{
	// 	neighbourfoodamteaten = curr - last
	// }

	receivedMsg := a.ReceiveMessage()
	if receivedMsg != nil {
		receivedMsg.Visit(a)
	} else {
		a.Log("I got nothing")
	}
	//TODO: Define a threshold limit for other agents to respond to our sent message.
	direction := rand.Intn(1)
	a.SendingMessage(direction)
	// msg := messages.NewRequestLeaveFoodMessage(a.ID(), a.Floor(), 10) //need to change how much to request to leave
	// a.SendMessage(-1, msg)
	// a.AppendToMessageMemory(-1, msg)
	// a.Log("I sent a message", infra.Fields{"message": "RequestLeaveFood"})
	// Verbose explanation of calculating the intented food taken
	// trust_prop := 100 - int(a.globalTrust)
	// intversion := int(a.CurrPlatFood())
	// food_prop := intversion * trust_prop / 100
	// a.IntendedFoodTaken = food.FoodType(food_prop)
	a.IntendedFoodTaken = food.FoodType(int(int(a.CurrPlatFood()) * (100 - int(a.globalTrust)) / 100))

	a.lastFoodTaken, _ = a.TakeFood(a.IntendedFoodTaken)
	// MessageToSend +=1
	a.MessageToSend += rand.Intn(15)
}

func (a *CustomAgent4) HandleAskHP(msg messages.AskHPMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.TargetFloor(), a.HP())
	a.SendMessage(reply)
	a.Log("I received an askHP message from ", infra.Fields{"floor": msg.SenderFloor(), "hp": a.HP()})
}

func (a *CustomAgent4) HandleAskFoodTaken(msg messages.AskFoodTakenMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.TargetFloor(), int(a.lastFoodTaken))
	a.SendMessage(reply)
	a.Log("I received an askFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": a.lastFoodTaken})
}

func (a *CustomAgent4) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.TargetFloor(), int(a.IntendedFoodTaken))
	a.SendMessage(reply)
	a.Log("I received an askIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": a.IntendedFoodTaken})
}

func (a *CustomAgent4) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	//fmt.Printf(msg.ID())
	reply := msg.Reply(a.ID(), a.Floor(), msg.TargetFloor(), true) // TODO: Change for later dependent on circumstance
	a.SendMessage(reply)
	a.Log("I received a requestLeaveFood message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func (a *CustomAgent4) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.TargetFloor(), true) // TODO: Change for later dependent on circumstance
	a.SendMessage(reply)
	a.Log("I received a requestTakeFood message from ", infra.Fields{"floor": msg.SenderFloor()})
}

func remove(slice MessageMemory, s int) ([]messages.Message, []int) {
	return append(slice.messages[:s], slice.messages[s+1:]...), append(slice.direction[:s], slice.direction[s+1:]...)
}

func (a *CustomAgent4) HandleResponse(msg messages.BoolResponseMessage) {
	response := msg.Response() // TODO: Change for later dependent on circumstance
	if !response {
		a.globalTrust += a.globalTrustSubtract * a.coefficients[0] // TODO: adapt for other conditions
	} else { // Iterating through all messages in agent memory
		for i := 0; i < len(a.sentMessages.messages); i++ {
			if msg.RequestID() == a.sentMessages.messages[i].ID() {
				a.Log("Received a message ", infra.Fields{"sender_uuid": msg.ID(), "sentmessage_uuid": a.sentMessages.messages[i].ID()})

				sentMessage := a.sentMessages.messages[i]
				sentMessageDirection := a.sentMessages.direction[i]
				a.sentMessages.messages, a.sentMessages.direction = remove(a.sentMessages, i) //a.sentMessages.messages[:i]+ a.sentMessages.messages[i+1:]

				// fooType := reflect.TypeOf(sentMessage)
				// 	for j := 0; j < fooType.NumMethod(); j++ {
				// 		method := fooType.Method(j)
				// 		fmt.Println(method.Name)

				//12 is RequestLeaveFoodMessage.MessageType(), 13 is RequestTakeFoodMessage.MessageType()

				if sentMessage.MessageType() == messages.RequestLeaveFood && sentMessageDirection == 1 {
					reqMessage, ok := sentMessage.(messages.RequestMessage)
					if !ok {
						err := fmt.Errorf("RequestMessage type assertion failed")
						fmt.Println(err.Error())
					} else if food.FoodType(reqMessage.Request()) <= a.CurrPlatFood() {
						a.globalTrust += a.globalTrustAdd * a.coefficients[1]
					} else if food.FoodType(reqMessage.Request()) > a.CurrPlatFood() {
						a.globalTrust += a.globalTrustSubtract * a.coefficients[2]
					}
				}

				// if sentMessage.MessageType() == messages.RequestTakeFood && sentMessageDirection == -1 { /////////////////////////////////////////////////////////////////////////
				// 	reqMessage := sentMessage.(messages.RequestMessage)
				// 	if food.FoodType(reqMessage.Request()) <= a.CurrPlatFood() {
				// 		a.globalTrust += a.globalTrustAdd * a.coefficients[1]
				// 	} else if food.FoodType(reqMessage.Request()) > a.CurrPlatFood() {
				// 		a.globalTrust += a.globalTrustSubtract * a.coefficients[2]
				// 	}
				// }

			}
		}
	}
	//a.Log("I received a Response message from ", infra.Fields{"floor": msg.SenderFloor(), "response": response})
}

func (a *CustomAgent4) HandleStateFoodTaken(msg messages.StateFoodTakenMessage) {
	statement := msg.Statement()
	if food.FoodType(statement) > a.maxFoodLimit {
		a.globalTrust += a.globalTrustSubtract * a.coefficients[3]
	} else {
		a.globalTrust += a.globalTrustAdd * a.coefficients[3]
	}
	a.Log("I received a StateFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": statement})
}

func (a *CustomAgent4) HandleStateHP(msg messages.StateHPMessage) {
	statement := msg.Statement()
	a.globalTrust += a.globalTrustAdd * a.coefficients[0]
	a.Log("I received a StateHP message from ", infra.Fields{"floor": msg.SenderFloor(), "hp": statement})
}

func (a *CustomAgent4) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {
	statement := msg.Statement()
	if food.FoodType(statement) > a.maxFoodLimit {
		a.globalTrust += a.globalTrustSubtract * a.coefficients[3]
	} else {
		a.globalTrust += a.globalTrustAdd * a.coefficients[3]
	}
	a.Log("I received a StateIntendedFoodTaken message from ", infra.Fields{"floor": msg.SenderFloor(), "food": statement})
}

// func (a *CustomAgent4) Run() {
// 	a.Log("Reporting agent state", infra.Fields{"health": a.HP(), "floor": a.ID(), a.Floor()})

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
// 			if receivedMsg.food == a.currPlatFood() && receivedMsg.senderFloor - a.ID(), a.Floor() == -1 && a.CurrPlatFood() != -1{ // on the floor above you
// 				a.globalTrust+= a.globalTrustAdd //
// 				if a.globalTrust > 100.0{
// 					a.globalTrust = 100.0
// 				}

// 			} else if receivedMsg.food != a.currPlatFood() && receivedMsg.senderFloor - a.Floor() == 1 && a.CurrPlatFood() != -1{ // on the floor below you

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
