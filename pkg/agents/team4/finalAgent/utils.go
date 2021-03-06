package team4EvoAgent

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
)

/*------------------------REMOVE A PARTICULAR MESSAGE FROM OUR MESSAGE ARRAY------------------------*/

func remove(slice []messages.Message, s int) []messages.Message {
	if s == len(slice)-1 {
		return slice[:s]
	}
	return append(slice[:s], slice[s+1:]...)
}

/*------------------------RETURNS HEALTH STATUS INFORMATION------------------------*/

func getHealthStatus(healthInfo *health.HealthInfo, healthLevelSeparation int, currentHp int) int {
	if currentHp <= healthInfo.WeakLevel { //critical
		return 0
	}
	if currentHp <= healthInfo.WeakLevel+healthLevelSeparation { //weak
		return 1
	}
	if currentHp <= healthInfo.WeakLevel+2*healthLevelSeparation { //normal
		return 2
	}
	//strong
	return 3
}

/*------------------------UTILITIES FOR MATCHING SENT MESSAGES------------------------*/

func getMatchingSentMessage(a *CustomAgentEvo, resMsg messages.ResponseMessage) messages.Message {
	var retMsg messages.Message
	for _, sentMsg := range a.params.sentMessages { // Iterate through each sent message
		if resMsg.RequestID() == sentMsg.ID() { // matches request to relevant sentID
			retMsg = sentMsg
			break
		}
	}
	return retMsg
}

func removeMatchingSentMessage(a *CustomAgentEvo, resMsg messages.ResponseMessage) {
	for i, sentMsg := range a.params.sentMessages { // Iterate through each sent message
		if resMsg.RequestID() == sentMsg.ID() { // matches request to relevant sentID
			a.params.sentMessages = remove(a.params.sentMessages, i)
			return
		}
	}
}
