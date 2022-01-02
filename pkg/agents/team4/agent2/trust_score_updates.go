package agentTrust

import (
	"fmt"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

func (a *CustomAgent4) CheckForResponse(msg messages.BoolResponseMessage) {
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
				}
			}
		}

	} else {
		for j := 0; j < len(a.sentMessages.messages); j++ {
			if msg.RequestID() == a.sentMessages.messages[j].ID() {
				sentMsg := a.sentMessages.messages[j]
				if sentMsg.MessageType() == messages.RequestTakeFood && a.NeighbourFoodEaten() == -1 {
					a.AppendToMessageMemory(a.Floor()-msg.SenderFloor(), &msg, a.responseMessages)
				} else if sentMsg.MessageType() == messages.RequestLeaveFood && !a.PlatformOnFloor() {
					a.AppendToMessageMemory(a.Floor()-msg.SenderFloor(), &msg, a.responseMessages)
				}
			}
			break
		}
	}
}
