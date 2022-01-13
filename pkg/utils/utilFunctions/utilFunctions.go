package utilFunctions

import "github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"

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
