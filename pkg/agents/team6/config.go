package team6

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/SOMAS2021/SOMAS2021/pkg/agents"
)

type behaviour float64

// const (
//  altruist behaviour = iota
//  collectivist
//  selfish
//  narcissist
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
	//maximum behaviour score an agent can reach
	maxBehaviourThreshold behaviour
}

type CustomAgent6 struct {
	*agents.Base
	config team6Config
	//keep track of the lowest floor we've been to
	maxFloorGuess int
	currBehaviour behaviour
}

var maxBehaviourThreshold behaviour = 10.0

type behaviourParameterWeights []float64

func chooseInitialBehaviour() behaviour {
	rand.Seed(time.Now().UnixNano())
	return behaviour(rand.Float64()) * maxBehaviourThreshold
}

func New(baseAgent *agents.Base) (agents.Agent, error) {
	initialBehaviour := chooseInitialBehaviour()
	return &CustomAgent6{
		Base: baseAgent,
		config: team6Config{
			baseBehaviour:         initialBehaviour,
			stubbornness:          0.5,
			maxBehaviourSwing:     2,
			paramWeights:          behaviourParameterWeights{0.7, 0.3}, //ensure sum of weights = max behaviour enum
			lambda:                3.0,
			maxBehaviourThreshold: maxBehaviourThreshold,
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

	var b behaviour = 2

	log.Printf("debug behaviour: " + b.String())

}

type thresholdBehaviourPair struct {
	threshold behaviour
	bType     string
}

func (b behaviour) String() string {
	//strings := [...]string{"Altruist", "Collectivist", "Selfish", "Narcissist"}

	behaviourMap := [...]thresholdBehaviourPair{{2, "Altruist"}, {7, "Collectivist"}, {9, "Selfish"}, {10, "Narcissist"}}

	if b >= 0 {

		for _, v := range behaviourMap {
			if b <= v.threshold {
				return v.bType
			}
		}
	}
	// if b >= 0 && int(b) < len(strings) {
	//  return strings[int(b)]
	// }
	return fmt.Sprintf("UNKNOWN Behaviour '%v'", int(b))
}
