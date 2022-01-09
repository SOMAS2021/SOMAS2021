package team2

func InitStateSpace(dim1 int, dim2 int, dim3 int, dim4 int) [][][][]int {
	//dim1: hp: 0-9, 10-19, ..., >90
	//dim2: foodOnPlatform 0-9, 10-19, ..., >90
	//dim3: daysAtCritical 0-8
	//dim4: neighbourHP: unknown = -1, 0-9, ... , >90
	stateSpace := make([][][][]int, dim1)
	stateNum := 0
	for i := 0; i < dim1; i++ {
		stateSpace[i] = make([][][]int, dim2)
		for j := 0; j < dim2; j++ {
			stateSpace[i][j] = make([][]int, dim3)
			for k := 0; k < dim3; k++ {
				stateSpace[i][j][k] = make([]int, dim4)
				for l := 0; l < dim4; l++ {
					stateSpace[i][j][k][l] = stateNum
					stateNum++
				}
			}
		}
	}
	return stateSpace
}

func (a *CustomAgent2) CheckState() int {
	return a.stateSpace[CategoriseObs(a.HP())][CategoriseObs(int(a.CurrPlatFood()))][a.DaysAtCritical()][CategoriseObs(a.neiboughHP)]
}

func CategoriseObs(observation int) (category int) {
	if observation == -1 {
		return 10
	}
	if observation >= 90 {
		observation = 90
	}
	return observation / 10
}
