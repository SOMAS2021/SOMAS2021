package team6

import (
	"fmt"
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
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
	maxFloorGuess      int
	currBehaviour      behaviour
	foodTakeDay        int
	reqLeaveFoodAmount int
	lastFoodTaken      food.FoodType
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
		currBehaviour:      initialBehaviour,
		maxFloorGuess:      baseAgent.Floor() + 2,
		foodTakeDay:        0,
		reqLeaveFoodAmount: -1,
		lastFoodTaken:      0,
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

	// Sending messages
	a.RequestLeaveFood()

	// Receiving messages
	receivedMsg := a.ReceiveMessage()
	if receivedMsg != nil {
		receivedMsg.Visit(a)
	} else {
		a.Log("I got no thing")
	}

	// a.Log("Custom agent 6 after update:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "behaviour": a.currBehaviour.String(), "maxFloorGuess": a.maxFloorGuess})

	foodTaken, err := a.TakeFood(a.intendedFoodIntake())
	if err != nil {
		switch err.(type) {
		case *infra.FloorError:
		case *infra.NegFoodError:
		case *infra.AlreadyEatenError:
		default:
		}
	} else {
		a.lastFoodTaken = foodTaken
	}

	a.Log("Team 6 took:", infra.Fields{"foodTaken": foodTaken, "bType": a.currBehaviour.String()})
	a.Log("Team 6 agent has HP:", infra.Fields{"hp": a.HP()})

	//fmt.Println(a.ActiveTreaties())

	// treaty := messages.NewTreaty(1, 1, 1, 1, 5, a.ID())
	// treatyMsg := messages.NewProposalMessage(a.ID(), a.Floor()+1, *treaty)

	// treatyMsg.Visit(a).

	treaty := messages.NewTreaty(1, 1, 1, 0, 1, 0, 5, a.ID())

	treaty.SignTreaty()
	a.AddTreaty(*treaty)

	treaty1 := messages.NewTreaty(1, 1, 1, 2, 1, 0, 5, a.ID())
	treaty1.SignTreaty()
	a.AddTreaty(*treaty1)
	// treaty2 := messages.NewTreaty(1, 1, 1, 10, 1, 4, 5, a.ID())
	// treaty2.SignTreaty()
	// a.AddTreaty(*treaty2)
	// treaty3 := messages.NewTreaty(1, 1, 1, 4, 1, 1, 5, a.ID())
	// treaty3.SignTreaty()
	// a.AddTreaty(*treaty3)
	// treaty4 := messages.NewTreaty(1, 1, 1, 11, 1, 3, 5, a.ID())
	// treaty4.SignTreaty()
	// a.AddTreaty(*treaty4)

	//treaty5 := messages.NewTreaty(1, 1, 1, 4, 1, 4, 5, a.ID())
	treaty5 := messages.NewTreaty(1, 1, 1, 101, 1, 1, 5, a.ID())

	c, d, e := a.foodRange(*treaty5)
	fmt.Println(c)
	fmt.Println(d)
	fmt.Println(e)

	// fmt.Println(a.ActiveTreaties())
	// fmt.Println(len(a.ActiveTreaties()))

	b := 0
	for b < 100 {

	}
}
