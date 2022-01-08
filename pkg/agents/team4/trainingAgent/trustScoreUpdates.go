package team4TrainingEvoAgent

import (
	"fmt"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

// Checks if neighbour below has eaten
func (a *CustomAgentEvo) neighbourFoodEaten() food.FoodType {
	if a.CurrPlatFood() != -1 {
		if !a.PlatformOnFloor() && a.CurrPlatFood() != a.params.lastPlatFood {
			return a.params.lastPlatFood - a.CurrPlatFood()
		}
		return 0
	}
	return -1
}

// removes a specific message from a message array
func remove(slice []messages.Message, s int) []messages.Message {
	return append(slice[:s], slice[s+1:]...)
}

func (a *CustomAgentEvo) CheckForResponse(msg messages.BoolResponseMessage) {
	if a.PlatformOnFloor() && len(a.params.responseMessages) > 0 { // Check if there are any responses messages.
		for i := 0; i < len(a.params.responseMessages); i++ { // Iterate through each response message
			respMsg := a.params.responseMessages[i]
			resMsg, ok := respMsg.(messages.ResponseMessage)
			if !ok {
				err := fmt.Errorf("ResponseMessage type assertion failed")
				fmt.Println(err.Error())
			} else {
				for j := 0; j < len(a.params.sentMessages); j++ { // Iterate through each sent message
					sentMsg := a.params.sentMessages[j]

					if resMsg.RequestID() == sentMsg.ID() { // Find the corresponding response message that's been sent
						a.params.sentMessages = remove(a.params.sentMessages, j) // Remove the accessed response/sent messages from memory
						a.params.responseMessages = remove(a.params.responseMessages, i)
						a.Log("Team4 received a message", infra.Fields{"sender_uuid": msg.ID(), "sentmessage_uuid": sentMsg.ID()})

						if sentMsg.MessageType() == messages.RequestLeaveFood && a.Floor()-msg.SenderFloor() == 1 { //TODO: theres now target floors and not directions anymore
							a.Log("Team4 reponse message received", infra.Fields{"sentMsg_Type": sentMsg.MessageType()})
							reqMsg, ok := sentMsg.(messages.RequestMessage)
							if !ok {
								err := fmt.Errorf("RequestMessage type assertion failed")
								fmt.Println(err.Error())
							} else if food.FoodType(reqMsg.Request()) <= a.CurrPlatFood() {
								a.AddToGlobalTrust(a.params.trustCoefficients[1])
							} else if food.FoodType(reqMsg.Request()) > a.CurrPlatFood() {
								a.SubFromGlobalTrust(a.params.trustCoefficients[1])
								a.Log("Team4: For Requested Food to Leave greater than Food on platform", infra.Fields{"Request_amt": reqMsg.Request(), "Food_on_our_level": a.CurrPlatFood(), "global_trust": a.params.globalTrust})
							}
						}
						break
					}
				}

			}

		}
	} else if a.params.lastFoodTaken+a.CurrPlatFood() != a.params.lastPlatFood {
		for i := 0; i < len(a.params.responseMessages); i++ { // Iterate through each response message
			respMsg := a.params.responseMessages[i]
			resMsg, ok := respMsg.(messages.ResponseMessage)
			if !ok {
				err := fmt.Errorf("ResponseMessage type assertion failed")
				fmt.Println(err.Error())
			} else {
				for j := 0; j < len(a.params.sentMessages); j++ { // Iterate through each sent message
					sentMsg := a.params.sentMessages[j]

					if resMsg.RequestID() == sentMsg.ID() { // Find the corresponding response message that's been sent
						a.params.sentMessages = remove(a.params.sentMessages, j) // Remove the accessed response/sent messages from memory
						a.params.responseMessages = remove(a.params.responseMessages, i)
						a.Log("Team4 received a message", infra.Fields{"sender_uuid": msg.ID(), "sentmessage_uuid": sentMsg.ID()})

						if sentMsg.MessageType() == messages.RequestTakeFood && a.Floor()-msg.SenderFloor() == -1 {
							reqMsg, ok := sentMsg.(messages.RequestMessage)
							if !ok {
								err := fmt.Errorf("RequestMessage type assertion failed")
								fmt.Println(err.Error())
							} else if food.FoodType(reqMsg.Request()) >= a.neighbourFoodEaten() {
								a.Log("Team4: For Requested Food to Take greater then or equal neighbour food eaten", infra.Fields{"Request_amt": reqMsg.Request(), "Food_on_our_level": a.neighbourFoodEaten(), "global_trust": a.params.globalTrust})
								a.AddToGlobalTrust(a.params.trustCoefficients[1])
							} else if food.FoodType(reqMsg.Request()) < a.neighbourFoodEaten() {
								a.Log("Team4: For Requested Food to Take less than neighbour food eaten", infra.Fields{"Request_amt": reqMsg.Request(), "Food_on_our_level": a.neighbourFoodEaten(), "global_trust": a.params.globalTrust})
								a.SubFromGlobalTrust(a.params.trustCoefficients[1])
							}
						}
						break
					}
				}
			}
		}

	} else {
		for j := 0; j < len(a.params.sentMessages); j++ {
			if msg.RequestID() == a.params.sentMessages[j].ID() {
				sentMsg := a.params.sentMessages[j]
				a.Log("Team4 received a message", infra.Fields{"sender_uuid": msg.ID(), "sentmessage_uuid": sentMsg.ID()})

				if sentMsg.MessageType() == messages.RequestTakeFood && a.neighbourFoodEaten() == -1 {
					a.params.sentMessages = append(a.params.responseMessages, &msg)
				} else if sentMsg.MessageType() == messages.RequestLeaveFood && !a.PlatformOnFloor() {
					a.params.sentMessages = append(a.params.responseMessages, &msg)
				} else if sentMsg.MessageType() == messages.RequestTakeFood && a.Floor()-msg.SenderFloor() == 1 {
					a.AddToGlobalTrust(a.params.trustCoefficients[1])
				} else if sentMsg.MessageType() == messages.RequestLeaveFood && a.Floor()-msg.SenderFloor() == -1 {
					a.AddToGlobalTrust(a.params.trustCoefficients[1])
				}
				break
			}
		}
	}
}

func (a *CustomAgentEvo) AddToGlobalTrust(coeff float64) {
	a.params.globalTrust += coeff
	if a.params.globalTrust > 100 {
		a.params.globalTrust = 100
	}
}

func (a *CustomAgentEvo) SubFromGlobalTrust(coeff float64) {
	a.params.globalTrust -= coeff
	if a.params.globalTrust < 0 {
		a.params.globalTrust = 0
	}

}
