package team4EvoAgent

import (
	"fmt"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
	log "github.com/sirupsen/logrus"
)

/*------------------------UPDATE GLOBAL TRUST UTLITY FUNCTION ------------------------*/

func (a *CustomAgentEvo) addToGlobalTrust(coeff float64) {
	a.params.globalTrust += coeff

	// Limiting global trust bounds between 0 and 100
	if a.params.globalTrust > 100 {
		a.params.globalTrust = 100
	} else if a.params.globalTrust < 0 {
		a.params.globalTrust = 0
	}
}

/*------------------------CHECK IF NEIGHBOUR BELOW HAS EATEN ------------------------*/

func (a *CustomAgentEvo) neighbourFoodEaten() food.FoodType {
	if a.CurrPlatFood() != -1 {
		if !a.PlatformOnFloor() && a.CurrPlatFood() != a.params.lastPlatFood {
			return a.params.lastPlatFood - a.CurrPlatFood()
		}
		return 0
	}
	return -1
}

/*------------------------TYPE ASSERTION UTILITIES FUNCTIONS------------------------*/

func (a *CustomAgentEvo) typeAssertResponseMessage(responseMessage messages.Message) (messages.ResponseMessage, error) {
	typeAsserted, ok := responseMessage.(messages.ResponseMessage)
	if !ok {
		err := fmt.Errorf("ResponseMessage type assertion failed")
		return nil, err
	}
	return typeAsserted, nil
}

func (a *CustomAgentEvo) typeAssertRequestMessage(requestMessage messages.Message) (messages.RequestMessage, error) {
	typeAsserted, ok := requestMessage.(messages.RequestMessage)
	if !ok {
		err := fmt.Errorf("RequestMessage type assertion failed")
		return nil, err
	}
	return typeAsserted, nil
}

/*------------------------HANDLING GLOBAL TRUST FROM MESSAGES ------------------------*/

func (a *CustomAgentEvo) updateGlobalTrustReqLeaveFood(msg messages.ResponseMessage, sentMsg messages.Message) {
	a.Log("Team4 reponse message received", infra.Fields{"sentMsg_Type": sentMsg.MessageType().String(), "senderFloor": msg.SenderFloor()})
	reqMsg, err := a.typeAssertRequestMessage(sentMsg)
	if err != nil {
		log.Error(err)
		return
	}
	if food.FoodType(reqMsg.Request()) <= a.CurrPlatFood() {
		a.addToGlobalTrust(a.params.trustCoefficients[1])
		a.Log("Team4: For Requested Food to Leave less than or equal to Food on platform", infra.Fields{"Request_amt": reqMsg.Request(), "Food_on_our_level": a.CurrPlatFood(), "global_trust": a.params.globalTrust})
		return
	} else {
		a.addToGlobalTrust(-a.params.trustCoefficients[1])
		a.Log("Team4: For Requested Food to Leave greater than Food on platform", infra.Fields{"Request_amt": reqMsg.Request(), "Food_on_our_level": a.CurrPlatFood(), "global_trust": a.params.globalTrust})
	}
}

func (a *CustomAgentEvo) updateGlobalTrustReqTakeFood(msg messages.ResponseMessage, sentMsg messages.Message) {
	a.Log("Team4 reponse message received", infra.Fields{"sentMsg_Type": sentMsg.MessageType().String(), "senderFloor": msg.SenderFloor()})
	reqMsg, err := a.typeAssertRequestMessage(sentMsg)
	if err != nil {
		log.Error(err)
		return
	}
	if food.FoodType(reqMsg.Request()) >= a.neighbourFoodEaten() {
		a.Log("Team4: For Requested Food to Take greater then or equal neighbour food eaten", infra.Fields{"Request_amt": reqMsg.Request(), "Food_on_our_level": a.neighbourFoodEaten(), "global_trust": a.params.globalTrust})
		a.addToGlobalTrust(a.params.trustCoefficients[1])
		return
	} else {
		a.Log("Team4: For Requested Food to Take less than neighbour food eaten", infra.Fields{"Request_amt": reqMsg.Request(), "Food_on_our_level": a.neighbourFoodEaten(), "global_trust": a.params.globalTrust})
		a.addToGlobalTrust(-a.params.trustCoefficients[1])
	}
}

/*------------------------HANDLING TRUST ON RESPONSES ------------------------*/

func (a *CustomAgentEvo) verifyResponses() {
	if len(a.params.responseMessages) > 0 {
		for i, respMsg := range a.params.responseMessages { // Iterate through each response message
			resMsg, err := a.typeAssertResponseMessage(respMsg)
			if err != nil {
				log.Error(err)
			} else {
				sentMsg := getMatchingSentMessage(a, resMsg)
				isHandled := false
				if a.PlatformOnFloor() && sentMsg.MessageType() == messages.RequestLeaveFood && a.Floor()-resMsg.SenderFloor() == 1 { // Check if there are any responses messages.
					a.updateGlobalTrustReqLeaveFood(resMsg, sentMsg)
					isHandled = true
				} else if a.params.lastFoodTaken+a.CurrPlatFood() != a.params.lastPlatFood && sentMsg.MessageType() == messages.RequestTakeFood && a.Floor()-resMsg.SenderFloor() == -1 {
					a.updateGlobalTrustReqTakeFood(resMsg, sentMsg)
					isHandled = true
				}
				if isHandled {
					removeMatchingSentMessage(a, resMsg)
					a.params.responseMessages = remove(a.params.responseMessages, i)
				}
			}
		}
		return
	}
}

func (a *CustomAgentEvo) checkAgentResponseMemory() {

	for id, item := range a.params.agentResponseMemory {
		if item[0] != -1 && item[1] != -1 {
			healthLevelSeparation := int(0.33 * float64(a.HealthInfo().MaxHP-a.HealthInfo().WeakLevel))
			HPstatus := getHealthStatus(a.HealthInfo(), healthLevelSeparation, item[0])
			if item[1] > a.params.foodToEat["selfish"][HPstatus] {
				a.addToGlobalTrust(-a.params.trustCoefficients[1])
			} else if item[1] <= a.params.foodToEat["selfless"][HPstatus] {
				a.addToGlobalTrust(a.params.trustCoefficients[1])
			}
			delete(a.params.agentResponseMemory, id)
		}
	}
}

func (a *CustomAgentEvo) checkForResponse(msg messages.BoolResponseMessage) {
	for _, sentMsg := range a.params.sentMessages {
		if msg.RequestID() == sentMsg.ID() {
			a.Log("Team4 received a message", infra.Fields{"sender_uuid": msg.ID(), "sentmessage_uuid": sentMsg.ID()})

			if sentMsg.MessageType() == messages.RequestTakeFood && a.Floor()-msg.SenderFloor() == 1 { //assuming honest
				a.addToGlobalTrust(a.params.trustCoefficients[1])
			} else if sentMsg.MessageType() == messages.RequestLeaveFood && a.Floor()-msg.SenderFloor() == -1 { //assuming honest
				a.addToGlobalTrust(a.params.trustCoefficients[1])
			} else {
				a.params.responseMessages = append(a.params.responseMessages, &msg)
			}
		}
	}
}
