package team6

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
)

func (a *CustomAgent6) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage) {
	healthInfo := a.HealthInfo()
	currentHP := a.HP()

	levels := levelsData{
		strongLevel:  healthInfo.MaxHP * 3 / 5,
		healthyLevel: healthInfo.MaxHP * 3 / 10,
		weakLevel:    healthInfo.MaxHP * 1 / 10,
		critLevel:    0,
	}

	reply := true
	switch a.currBehaviour.String() {
	case "Altruist":
		reply = true

	case "Collectivist":
		if currentHP >= levels.weakLevel {
			reply = true
		} else {
			reply = false
		}

	case "Selfish":
		if currentHP >= levels.strongLevel {
			reply = true
		} else {
			reply = false
		}

	case "Narcissist":
		reply = false

	default:
		reply = true
	}

	replyMessage := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), reply)
	a.SendMessage(replyMessage)

	if reply == true {
		a.reqLeaveFoodAmount = msg.Request()
		a.Log("I received a requestTakeFood message and my response was true")
	} else {
		a.reqLeaveFoodAmount = -1
		a.Log("I received a requestTakeFood message and my response was false")
	}

}
