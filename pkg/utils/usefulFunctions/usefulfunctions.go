package usefulfunctions

func Sum(inputList []int) int {
	totalAgents := 0
	for _, value := range inputList {
		totalAgents += value
	}
	return totalAgents
}
