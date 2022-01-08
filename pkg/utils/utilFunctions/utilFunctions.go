package utilFunctions

import (
	"math"

	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
)

func Sum(input map[agent.AgentType]int) int {
	totalAgents := 0
	for _, value := range input {
		totalAgents += value
	}
	return totalAgents

}

func MinInt(vars ...int) int {
	min := vars[0]
	for _, i := range vars {
		if min > i {
			min = i
		}
	}
	return min
}

/*
Restricts 'value' to [lowerBound, uppwerBound].
Values at the either bound will take on the value of that respective bound.
For single sided bounds, use +ve or -ve inf appropriately.
Returns +ve infinity on bounds error.
*/
func RestrictToRange(lowerBound, upperBound, value float64) float64 {
	if lowerBound >= upperBound {
		return math.Inf(1)
	}
	if value < lowerBound {
		return lowerBound
	}
	if value > upperBound {
		return upperBound
	}
	return value
}
