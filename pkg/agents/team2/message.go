package team2

import "github.com/SOMAS2021/SOMAS2021/pkg/messages"

func (a *CustomAgent2) AskNeighbourHP() {
	msg := messages.NewAskHPMessage(a.ID(), a.Floor(), a.Floor()+1)
	a.SendMessage(msg)
}

func (a *CustomAgent2) CheckNeighbourHP() {
	msg := a.ReceiveMessage()
	for msg != nil {
		if msg.MessageType() == 7 {
			msg.Visit(a)
		}
		msg = a.ReceiveMessage()
	}
}

func (a *CustomAgent2) ReplyAllAskMsg() {
	msg := a.ReceiveMessage()
	for msg != nil {
		msgType := msg.MessageType()
		if msgType == 1 || msgType == 2 || msgType == 4 {
			msg.Visit(a)
		}
		msg = a.ReceiveMessage()
	}
}
