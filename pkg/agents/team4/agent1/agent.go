package team4EvoAgent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

// type Coefficient struct {
// 	Floor []float64
// 	Hp    []float64
// }

// type Equation struct {
// 	coefficients []float64
// }

// // function to evaluate a function f(x) at a specific value of x0
// func (equation *Equation) EvaluateEquation(input int) float64 {
// 	sum := 0.0
// 	for i, coeff := range equation.coefficients {
// 		power := math.Pow(float64(input), float64(i))
// 		sum += coeff * power
// 	}
// 	return sum
// }

// // function to pick the correct co-efficients from the configuration file
// func GenerateEquations() (Equation, Equation) {
// 	mydir, err := os.Getwd()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	file, _ := ioutil.ReadFile(fmt.Sprintf("%s/pkg/agents/team4/agent1/agentConfig.json", mydir))
// 	var data1 Coefficient
// 	_ = json.Unmarshal(file, &data1) //parse the config json file to get the coeffs for floor and HP equations

// 	var newAgentFloor Equation
// 	var newAgentHp Equation

// 	newAgentFloor.coefficients = data1.Floor
// 	newAgentHp.coefficients = data1.Hp

// 	return newAgentFloor, newAgentHp
// }

type CustomAgentEvoParams struct {
	// these params are based only on singular inputs
	// scalingEquation   Equation
	// currentFloorScore Equation
	// currentHpScore    Equation
	foodToEat    []int
	daysToWait   []int
	ageLastEaten int
	// below params updated based in previous experience of floors
	trustscore float64
	morality   float64
	trauma     float64
}

type CustomAgentEvo struct {
	*infra.Base
	// new params
	params CustomAgentEvoParams
}

type LoadedData struct {
	FoodToEat  []int
	DaysToWait []int
}

func InitaliseParams(baseAgent *infra.Base) CustomAgentEvoParams {
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	file, _ := ioutil.ReadFile(fmt.Sprintf("%s/pkg/agents/team4/agent1/agentConfig.json", mydir))
	var data1 LoadedData
	_ = json.Unmarshal(file, &data1) //parse the config json file to get the coeffs for floor and HP equations

	return CustomAgentEvoParams{ //initialise the parameters of the agent
		// scalingEquation: Equation{
		// 	coefficients: []float64{1, 1, 1, 1},
		// },
		// currentFloorScore: currentFloorScoreEquation,
		// currentHpScore:    currentHpScoreEquation,
		foodToEat:    data1.FoodToEat,
		daysToWait:   data1.DaysToWait,
		ageLastEaten: 0,
		trustscore:   100 * rand.Float64(), //random for now
		morality:     100 * rand.Float64(), //random for now
		trauma:       0,                    //0 for now
	}
}

func New(baseAgent *infra.Base) (infra.Agent, error) {
	// currentFloorScoreEquation, currentHpScoreEquation := GenerateEquations()
	//create other parameters
	return &CustomAgentEvo{
		Base:   baseAgent,
		params: InitaliseParams(baseAgent),
	}, nil
}

func (a *CustomAgentEvo) Run() {

	var healthStatus int
	var healthLevelSeparation = int(0.33 * float64(a.HealthInfo().MaxHP-a.HealthInfo().WeakLevel))

	if a.HP() <= a.HealthInfo().WeakLevel { //critical
		healthStatus = 0
	} else if a.HP() <= a.HealthInfo().WeakLevel+healthLevelSeparation { //weak
		healthStatus = 1
	} else if a.HP() <= a.HealthInfo().WeakLevel+2*healthLevelSeparation { //normal
		healthStatus = 2
	} else { //strong
		healthStatus = 3
	}

	var foodEaten food.FoodType
	var err error
	if (a.Age() - a.params.ageLastEaten) >= a.params.daysToWait[healthStatus] {
		foodEaten, err = a.TakeFood(food.FoodType(int(a.params.foodToEat[healthStatus])))
		a.params.ageLastEaten = a.Age()
		if err != nil {
			switch err.(type) {
			case *infra.FloorError:
			case *infra.NegFoodError:
			case *infra.AlreadyEatenError:
			default:
			}
		}
	}

	// // Using custom params to decide how much food to eat
	// scaledFloorScore := a.params.currentFloorScore.EvaluateEquation(a.Floor()) / a.params.scalingEquation.EvaluateEquation(a.Floor())
	// scaledHpScore := a.params.currentHpScore.EvaluateEquation(a.HP()) / a.params.scalingEquation.EvaluateEquation(a.HP())
	// foodToEat := food.FoodType(math.Max(0.0, 50*scaledFloorScore+50*scaledHpScore))
	// foodEaten := a.TakeFood(foodToEat)
	a.Log("team4EvoAgent reporting status:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "FoodToEat": a.params.foodToEat[healthStatus], "DaysToWait": a.params.daysToWait, "foodEaten": foodEaten})

	// // Trauma score for future use
	// if a.HP() < 25 {
	// 	a.params.trauma = math.Min(100, a.params.trauma+1)
	// } else if a.HP() > 75 {
	// 	a.params.trauma = math.Max(0, a.params.trauma-1)
	// }

}
