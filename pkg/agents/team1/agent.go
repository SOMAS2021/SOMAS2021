package agent

import (
	"log"

	baseagent "github.com/SOMAS2021/SOMAS2021/pkg/agents/default"
)

type CustomAgent struct {
	baseagent.BaseAgent
}

func (a *CustomAgent) Run() {
	log.Printf("Custom agent has floor: %d", a.GetFloor())
}
