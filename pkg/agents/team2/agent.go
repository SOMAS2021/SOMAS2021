package agent

import (
	"log"

	baseagent "github.com/SOMAS2021/SOMAS2021/pkg/agents/default"
)

type CustomAgent2 struct {
	// these should be capital to enable access from the tower
	BA       *baseagent.BaseAgent
	myNumber int
}

func New(baseAgent *baseagent.BaseAgent) (*CustomAgent2, error) {
	return &CustomAgent2{
		BA:       baseAgent,
		myNumber: 0,
	}, nil
}

func (a *CustomAgent2) Run() {
	log.Printf("Custom agent team 2 has floor: %d", a.BA.GetFloor())
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
