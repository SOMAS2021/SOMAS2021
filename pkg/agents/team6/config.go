package team6

import (
	"fmt"
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
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
	*infra.Base
	config team6Config
	//keep track of the lowest floor we've been to
	maxFloorGuess int
	currBehaviour behaviour
	foodTakeDay   int
}

type thresholdBehaviourPair struct {
	threshold behaviour
	bType     string
}

type behaviourParameterWeights []float64

var maxBehaviourThreshold behaviour = 10.0

func chooseInitialBehaviour() behaviour {
	return behaviour(rand.Float64()) * maxBehaviourThreshold
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	initialBehaviour := chooseInitialBehaviour()
	return &CustomAgent6{
		Base: baseAgent,
		config: team6Config{
			baseBehaviour:         initialBehaviour,
			stubbornness:          0.0,
			maxBehaviourSwing:     8,
			paramWeights:          behaviourParameterWeights{0.7, 0.3}, //ensure sum of weights = max behaviour enum
			lambda:                3.0,
			maxBehaviourThreshold: maxBehaviourThreshold,
		},
		currBehaviour: initialBehaviour,
		maxFloorGuess: baseAgent.Floor() + 2,
		foodTakeDay:   0,
	}, nil
}

func (b behaviour) String() string {
	behaviourMap := [...]thresholdBehaviourPair{{2, "Altruist"}, {7, "Collectivist"}, {9, "Selfish"}, {10, "Narcissist"}}

	if b >= 0 {
		for _, v := range behaviourMap {
			if b <= v.threshold {
				return v.bType
			}
		}
	}

	return fmt.Sprintf("UNKNOWN Behaviour '%v'", int(b))
}

func (a *CustomAgent6) Run() {
	// a.Log("Custom agent 6 before update:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "behaviour": a.currBehaviour.String(), "maxFloorGuess": a.maxFloorGuess})

	a.updateBehaviour()

	// a.Log("Custom agent 6 after update:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "behaviour": a.currBehaviour.String(), "maxFloorGuess": a.maxFloorGuess})

	foodAmount := a.foodIntake()
	a.TakeFood(foodAmount)
	a.Log("Team 6 took:", infra.Fields{"foodTaken": foodAmount, "bType": a.currBehaviour.String()})
	a.Log("Team 6 agent has HP:", infra.Fields{"hp": a.HP()})

	msg := messages.NewAskHPMessage(a.ID(), a.Floor())
	a.SendMessage(1, msg)
	a.Log("Team 6 sent message:", infra.Fields{"floor": a.Floor(), "messageType": msg.MessageType().String()})

}
