package team2

func InitStateSpace(dim1 int, dim2 int, dim3 int) [][][]int {
	stateSpace := make([][][]int, dim1)
	stateNum := 0
	for i := 0; i < dim1; i++ {
		stateSpace[i] = make([][]int, dim2)
		for j := 0; j < 3; j++ {
			stateSpace[i][j] = make([]int, dim3)
			for k := 0; k < 3; k++ {
				stateSpace[i][j][k] = stateNum
				stateNum++
			}
		}
	}
	return stateSpace
}

//TODO: the CheckState function can be replaced later
func (a *CustomAgent2) CheckState() int {
	floorState := CategoriseObs(a.Floor())
	if a.CurrPlatFood() == -1.0 {
		return -1
	}
	return a.stateSpace[a.HP()][floorState][CategoriseObs(int(a.CurrPlatFood()))]
}

func CategoriseObs(observation int) (category int) {
	if observation >= 61 {
		return 0
	} else if observation >= 31 {
		return 1
	} else {
		return 2
	}
}
