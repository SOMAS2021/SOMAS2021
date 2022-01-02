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
	globalTrust           float32
	globalTrustAdd        float32
	globalTrustSubtract   float32
	coefficients          []float32
	lastFoodTaken         food.FoodType
	IntendedFoodTaken     food.FoodType
	sentMessages          MessageMemory
	responseMessages      MessageMemory
	MessageToSend         int
	lastPlatFood          food.FoodType
	maxFoodLimit          food.FoodType
	neighbourFoodEatenAmt food.FoodType
}

type MessageMemory struct {
	// An array of messages stored into the agent before platform change.
	direction []int
	messages  []messages.Message
}

func (a *CustomAgent4) AppendToMessageMemory(direction int, msg messages.Message, msgMemory MessageMemory) {
	msgMemory.direction = append(msgMemory.direction, direction)
	msgMemory.messages = append(msgMemory.messages, msg)
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
	a.AppendToMessageMemory(direction, msg, a.sentMessages)
	a.Log("I sent a message", infra.Fields{"message": msg.MessageType()})

}

func (a *CustomAgent4) NeighbourFoodEaten() food.FoodType {
	if a.CurrPlatFood() != -1 {
		if !a.PlatformOnFloor() && a.CurrPlatFood() != a.lastPlatFood {
			return a.lastPlatFood - a.CurrPlatFood()
		}
		return 0
	}
	return -1
}

func remove(slice MessageMemory, s int) ([]messages.Message, []int) {
	return append(slice.messages[:s], slice.messages[s+1:]...), append(slice.direction[:s], slice.direction[s+1:]...)
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
		responseMessages: MessageMemory{
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

	a.IntendedFoodTaken = food.FoodType(int(int(a.CurrPlatFood()) * (100 - int(a.globalTrust)) / 100))

	a.lastFoodTaken, _ = a.TakeFood(a.IntendedFoodTaken)
	// MessageToSend +=1
	a.MessageToSend += rand.Intn(15)
}

func (a *CustomAgent4) CheckForResponse(msg messages.BoolResponseMessage) {
	if !msg.Response() {
		a.globalTrust += a.globalTrustSubtract * a.coefficients[0] // TODO: adapt for other conditions
	} else { // Iterating through all messages in agent memory

		if a.PlatformOnFloor() && len(a.responseMessages.messages) > 0 { // Check if there are any responses messages.
			for i := 0; i < len(a.responseMessages.messages); i++ { // Iterate through each response message
				respMsg := a.responseMessages.messages[i]
				resMsg, ok := respMsg.(messages.ResponseMessage)
				if !ok {
					err := fmt.Errorf("ResponseMessage type assertion failed")
					fmt.Println(err.Error())
				} else {
					for j := 0; j < len(a.sentMessages.messages); j++ { // Iterate through each sent message
						sentMsg := a.sentMessages.messages[j]
						sentMsgDir := a.sentMessages.direction[j]

						if resMsg.RequestID() == sentMsg.ID() { // Find the corresponding response message that's been sent
							a.sentMessages.messages, a.sentMessages.direction = remove(a.sentMessages, j) // Remove the accessed response/sent messages from memory
							a.responseMessages.messages, a.responseMessages.direction = remove(a.responseMessages, i)
							a.Log("Received a message", infra.Fields{"sender_uuid": msg.ID(), "sentmessage_uuid": sentMsg.ID()})

							if sentMsg.MessageType() == messages.RequestLeaveFood && sentMsgDir == 1 { //TODO: theres now target floors and not directions anymore
								a.Log("Reponse message received", infra.Fields{"sentMsg_Type": sentMsg.MessageType()})
								reqMsg, ok := sentMsg.(messages.RequestMessage)
								if !ok {
									err := fmt.Errorf("RequestMessage type assertion failed")
									fmt.Println(err.Error())
								} else if food.FoodType(reqMsg.Request()) <= a.CurrPlatFood() {
									a.globalTrust += a.globalTrustAdd * a.coefficients[1]
									a.Log("For Requested Food to Leave less than or equal to Food on platform", infra.Fields{"Request_amt": reqMsg.Request(), "Food_on_our_level": a.CurrPlatFood(), "global_trust": a.globalTrust})
								} else if food.FoodType(reqMsg.Request()) > a.CurrPlatFood() {
									a.globalTrust += a.globalTrustSubtract * a.coefficients[2]
									a.Log("For Requested Food to Leave greater than Food on platform", infra.Fields{"Request_amt": reqMsg.Request(), "Food_on_our_level": a.CurrPlatFood(), "global_trust": a.globalTrust})

								}
							}
						}
					}
					break
				}

			}
		} else if a.lastFoodTaken+a.CurrPlatFood() != a.lastPlatFood {
			for i := 0; i < len(a.responseMessages.messages); i++ { // Iterate through each response message
				respMsg := a.responseMessages.messages[i]
				resMsg, ok := respMsg.(messages.ResponseMessage)
				if !ok {
					err := fmt.Errorf("ResponseMessage type assertion failed")
					fmt.Println(err.Error())
				} else {
					for j := 0; j < len(a.sentMessages.messages); j++ { // Iterate through each sent message
						sentMsg := a.sentMessages.messages[j]
						sentMsgDir := a.sentMessages.direction[j]

						if resMsg.RequestID() == sentMsg.ID() { // Find the corresponding response message that's been sent
							a.sentMessages.messages, a.sentMessages.direction = remove(a.sentMessages, j) // Remove the accessed response/sent messages from memory
							a.responseMessages.messages, a.responseMessages.direction = remove(a.responseMessages, i)
							a.Log("Received a message", infra.Fields{"sender_uuid": msg.ID(), "sentmessage_uuid": sentMsg.ID()})

							if sentMsg.MessageType() == messages.RequestTakeFood && sentMsgDir == -1 {
								reqMessage, ok := sentMsg.(messages.RequestMessage)
								if !ok {
									err := fmt.Errorf("RequestMessage type assertion failed")
									fmt.Println(err.Error())
								} else if food.FoodType(reqMessage.Request()) >= a.NeighbourFoodEaten() {
									a.globalTrust += a.globalTrustAdd * a.coefficients[1]
								} else if food.FoodType(reqMessage.Request()) <= a.NeighbourFoodEaten() {
									a.globalTrust += a.globalTrustSubtract * a.coefficients[2]
								}
							}
						}
						break
					}
				}
			}

		} else {
			for j := 0; j < len(a.sentMessages.messages); j++ {
				if msg.RequestID() == a.sentMessages.messages[j].ID() {
					sentMsg := a.sentMessages.messages[j]
					if sentMsg.MessageType() == messages.RequestTakeFood && a.NeighbourFoodEaten() == -1 {
						a.AppendToMessageMemory(a.Floor()-msg.SenderFloor(), msg, a.responseMessages)
					} else if sentMsg.MessageType() == messages.RequestLeaveFood && !a.PlatformOnFloor() {
						a.AppendToMessageMemory(a.Floor()-msg.SenderFloor(), msg, a.responseMessages)
					}
				}
				break
			}
		}
	}
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

func (a *CustomAgent4) HandleResponse(msg messages.BoolResponseMessage) {
	response := msg.Response() // TODO: Change for later dependent on circumstance
	a.CheckForResponse(msg)
	a.Log("I received a Response message from ", infra.Fields{"floor": msg.SenderFloor(), "response": response})
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
