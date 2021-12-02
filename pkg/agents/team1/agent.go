package agent

import (
	"log"
	messages "github.com/SOMAS2021/SOMAS2021/pkg/infra/messages"
	baseagent "github.com/SOMAS2021/SOMAS2021/pkg/agents/default"
)

type CustomAgent struct {
	baseagent.BaseAgent
}

func (a *CustomAgent) Run() {
	//make message
	msg := messages.Message{a.BaseAgent.GetFloor()}
	a.baseagent.tower.sendMessage(-1, a, msg)
	recieved_msg := a.baseagent.recieveMessage()
	log.Printf("Custom agent has floor: %d, I have recieved message from agent", recieved_msg.GetSenderID())

}
