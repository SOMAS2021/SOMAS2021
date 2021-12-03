package team6

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/agents"
)

type behaviour float64

// const (
// 	altruist behaviour = iota
// 	collectivist
// 	selfish
// 	narcissist
// )

type team6Config struct {
	baseBehaviour behaviour
	//the scaling factor which limits the change in agent behaviour
	stubbornness float64
	//the largest jump in behaviour an agent can take
	maxBehaviourSwing float64
	//weights used to assess score for behaviour update
	paramWeights behaviourParameterWeights
	//floor scaling discount factor
	lambda float64
}

type CustomAgent6 struct {
	*agents.Base
	config team6Config
	//keep track of the lowest floor we've been to
	maxFloorGuess int
	currBehaviour behaviour
}

type behaviourParameterWeights []float64

func chooseInitialBehaviour() behaviour {
	return behaviour(rand.Intn(4))
}

func New(baseAgent *agents.Base) (agents.Agent, error) {
	initialBehaviour := chooseInitialBehaviour()
	return &CustomAgent6{
		Base: baseAgent,
		config: team6Config{
			baseBehaviour:     initialBehaviour,
			stubbornness:      0.5,
			maxBehaviourSwing: 1,
			paramWeights:      behaviourParameterWeights{2.0, 10.0},
			lambda:            3.0,
		},
		currBehaviour: initialBehaviour,
		maxFloorGuess: baseAgent.Floor() + 2,
	}, nil
}

func (a *CustomAgent6) Run() {
	log.Printf("Custom agent team 6 has floor: %d", a.Floor())
	log.Printf("Team 6 has behaviour: " + a.currBehaviour.String())
	log.Printf("Team 6 has maxFloorGuess: %d", a.maxFloorGuess)
	a.updateBehaviour()
	log.Printf("Team 6 has behaviour: " + a.currBehaviour.String())
	log.Printf("Team 6 has maxFloorGuess: %d", a.maxFloorGuess)

}

func (b behaviour) String() string {
	strings := [...]string{"Altruist", "Collectivist", "Selfish", "Narcissist"}
	if b >= 0 && int(b) < len(strings) {
		return strings[int(b)]
	}
	return fmt.Sprintf("UNKNOWN EmotionalState '%v'", int(b))
}
