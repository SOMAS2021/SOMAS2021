package team4EvoAgent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

type Coefficient struct {
	Floor [][]float64
	Hp    [][]float64
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
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	file, _ := ioutil.ReadFile(fmt.Sprintf("%s/pkg/agents/team4/agent1/config.json", mydir))
	var data1 Coefficient
	_ = json.Unmarshal(file, &data1)

	floorCoeffArr := data1.Floor
	hpCoeffArr := data1.Hp

	var newAgentFloor Equation
	var newAgentHp Equation

	for _, coeffArr := range floorCoeffArr {
		randVal := rand.Float64()
		if randVal < 0.65 {
			newAgentFloor.coefficients = append(newAgentFloor.coefficients, coeffArr[0])
		} else if randVal < 0.85 {
			newAgentFloor.coefficients = append(newAgentFloor.coefficients, coeffArr[1])
		} else if randVal < 0.99 {
			newAgentFloor.coefficients = append(newAgentFloor.coefficients, coeffArr[2])
		} else {
			newAgentFloor.coefficients = append(newAgentFloor.coefficients, rand.Float64())
		}
	}

	for _, coeffArr := range hpCoeffArr {
		randVal := rand.Float64()
		if randVal < 0.65 {
			newAgentHp.coefficients = append(newAgentHp.coefficients, coeffArr[0])
		} else if randVal < 0.85 {
			newAgentHp.coefficients = append(newAgentHp.coefficients, coeffArr[1])
		} else if randVal < 0.99 {
			newAgentHp.coefficients = append(newAgentHp.coefficients, coeffArr[2])
		} else {
			newAgentHp.coefficients = append(newAgentHp.coefficients, rand.Float64())
		}
	}

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

	beforeHP := a.HP()
	a.TakeFood(food.FoodType(foodToEat))
	foodEaten := a.HP() - beforeHP
	a.Log("team4EvoAgent reporting status:", infra.Fields{"floor": a.Floor(), "hp": a.HP(), "foodToEat": foodToEat, "foodEaten": foodEaten, "currentFloorScore": a.params.currentFloorScore.coefficients, "currentHpScore": a.params.currentHpScore.coefficients})
	if a.HP() < 25 {
		a.params.trauma = math.Min(100, a.params.trauma+1)

	} else if a.HP() > 75 {
		a.params.trauma = math.Max(0, a.params.trauma-1)
	}
}
