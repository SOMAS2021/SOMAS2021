package utilfunctions

func Sum(inputList []int) int {
	totalAgents := 0
	for _, value := range inputList {
		totalAgents += value
	}
	return totalAgents

}

func MinInt(valueOne, valueTwo int) int {
	if valueOne > valueTwo {
		return valueOne
	}
	return valueTwo
}
