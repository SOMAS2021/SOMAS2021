package team4TrainingEvoAgent

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/health"
)

/*------------------------REMOVE A PARTICULAR MESSAGE FROM OUR MESSAGE ARRAY------------------------*/

func remove(slice []messages.Message, s int) []messages.Message {
	retSlice := []messages.Message{}
	for i, value := range slice{
		if i != s{
			retSlice = append(retSlice,value )
		}
	}
	return retSlice
}

/*------------------------RETURNS HEALTH STATUS INFORMATION------------------------*/

func getHealthStatus(healthInfo *health.HealthInfo, healthLevelSeparation int, currentHp int) int {
	if currentHp <= healthInfo.WeakLevel { //critical
		return 0
	} else if currentHp <= healthInfo.WeakLevel+healthLevelSeparation { //weak
		return 1
	} else if currentHp <= healthInfo.WeakLevel+2*healthLevelSeparation { //normal
		return 2
	} else { //strong
		return 3
	}
}

/*------------------------UTILITIES FOR MATCHING SENT MESSAGES------------------------*/

// func getMatchingSentMessage(a *CustomAgentEvo, resMsg messages.ResponseMessage) messages.Message {
// 	var retMsg messages.Message
// 	for _, sentMsg := range a.params.sentMessages { // Iterate through each sent message
// 		if resMsg.RequestID() == sentMsg.ID() { // matches request to relevant sentID
// 			retMsg = sentMsg
// 			break
// 		}
// 	}
// 	return retMsg
// }

// func removeMatchingSentMessage(a *CustomAgentEvo, resMsg messages.ResponseMessage) {
// 	for i, sentMsg := range a.params.sentMessages { // Iterate through each sent message
// 		if resMsg.RequestID() == sentMsg.ID() { // matches request to relevant sentID
// 			a.params.sentMessages = remove(a.params.sentMessages, i)
// 			break
// 		}
// 	}
// }
