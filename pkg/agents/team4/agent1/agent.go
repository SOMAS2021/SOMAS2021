package team4EvoAgent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"sync"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

type Coefficient struct {
	Floor []float64
	Hp    []float64
}

type Equation struct {
	coefficients []float64
	
}

func (equation *Equation) EvaluateEquation(input int) float64 {
	sum := 0.0
	for i, coeff := range equation.coefficients {
		power := math.Pow(float64(input), float64(i))
		sum += coeff * power
	}
	return sum
}

func GenerateEquations() (Equation, Equation) {
	var mu	sync.Mutex
	mu.Lock()
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	file, _ := ioutil.ReadFile(fmt.Sprintf("%s/pkg/agents/team4/agent1/agentConfig.json", mydir))
	var data1 Coefficient
	_ = json.Unmarshal(file, &data1)

	var newAgentFloor Equation
	var newAgentHp Equation

	newAgentFloor.coefficients = data1.Floor
	newAgentHp.coefficients = data1.Hp

	defer mu.Unlock()
	return newAgentFloor, newAgentHp
}

type CustomAgentEvoParams struct {
	// these params are based only on singular inputs
	currentFloorScore Equation
	currentHpScore    Equation
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

func New(baseAgent *infra.Base) (infra.Agent, error) {
	currentFloorScoreEquation, currentHpScoreEquation := GenerateEquations()
	//create other parameters
	return &CustomAgentEvo{
		Base: baseAgent,
		params: CustomAgentEvoParams{
			currentFloorScore: currentFloorScoreEquation,
			currentHpScore:    currentHpScoreEquation,
			trustscore:        100 * rand.Float64(),
			morality:          100 * rand.Float64(),
			trauma:            0,
		},
	}, nil
}

func (a *CustomAgentEvo) Run() {
	foodToEat := a.params.currentFloorScore.EvaluateEquation(a.Floor()) + a.params.currentHpScore.EvaluateEquation(a.HP())
	// scaledFoodToEat := 100 * math.Tanh(foodToEat/5000)
	scaledFoodToEat := foodToEat
	// fmt.Printf("%f", 100*math.Tanh(foodToEat/5000))

	beforeHP := a.HP()
	a.TakeFood(food.FoodType(foodToEat))
	foodEaten := a.HP() - beforeHP
	a.Log("team4EvoAgent reporting status:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "foodToEat": scaledFoodToEat, "foodEaten": foodEaten, "currentFloorScore": a.params.currentFloorScore.coefficients, "currentHpScore": a.params.currentHpScore.coefficients})
	// fmt.Printf("Food to eat: %f", foodToEat)
	if a.HP() < 25 {
		a.params.trauma = math.Min(100, a.params.trauma+1)
	} else if a.HP() > 75 {
		a.params.trauma = math.Max(0, a.params.trauma-1)
	}
}
