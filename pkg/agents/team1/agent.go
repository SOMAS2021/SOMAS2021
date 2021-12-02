package agent

import (
	"log"

	baseagent "github.com/SOMAS2021/SOMAS2021/pkg/agents/default"
)

type CustomAgent1 struct {
	// these should be capital to enable access from the tower
	BA       *baseagent.BaseAgent
	myString string
}

func New(baseAgent *baseagent.BaseAgent) (baseagent.Agent, error) {
	return &CustomAgent1{
		BA:       baseAgent,
		myString: "test",
	}, nil
}

func (a *CustomAgent1) Run() {
	log.Printf("Custom agent team 1 has floor: %d", a.BA.GetFloor())
}

// package agent

// import (
// 	"log"

// 	baseagent "github.com/SOMAS2021/SOMAS2021/pkg/agents/default"
// )

// type CustomAgent struct {
// 	// these should be capital to enable access from the tower
// 	BA *baseagent.BaseAgent
// }

// type Pointer struct {
// 	CA *CustomAgent
// }

// func New(baseAgent *baseagent.BaseAgent) (interface{}, error) {
// 	t := &CustomAgent{
// 		BA: baseAgent,
// 	}
// 	return Pointer{
// 		CA: t,
// 	}, nil
// }

// func (a *CustomAgent) Run() {
// 	log.Printf("Custom agent has floor: %d", a.BA.GetFloor())
// }
