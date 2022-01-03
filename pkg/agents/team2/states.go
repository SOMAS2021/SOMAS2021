package team2

func InitStateSpace(dim1 int, dim2 int, dim3 int, dim4 int) [][][][]int {
	//dim1: hp 0-100
	//dim2: floor
	//dim3: foodOnPlatform
	//dim2: daysAtCritical 0-8
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

//TODO: the CheckState function can be replaced later
func (a *CustomAgent2) CheckState() int {
	floorState := CategoriseObs(a.Floor())
	if a.CurrPlatFood() == -1.0 {
		return -1
	}
	return a.stateSpace[a.HP()][floorState][CategoriseObs(int(a.CurrPlatFood()))][a.DaysAtCritical()]
}

func CategoriseObs(observation int) (category int) {
	if float64(observation) >= 0.666*float64(observation) {
		return 0
	} else if float64(observation) >= 0.333*float64(observation) {
		return 1
	} else {
		return 2
	}
}
