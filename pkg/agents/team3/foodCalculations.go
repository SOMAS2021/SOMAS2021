package team3

import (
	"math"

	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

/*
stubbornness range = 0-100
mood range = 0-100
morality range = 0-100
*/

func takeFoodCalculation(a *CustomAgent3) food.FoodType {
	//food taken is: 10 - morality/20 - mood/20
	return food.FoodType(10.0 - math.Floor((float64(a.vars.morality) / 20.0)) - math.Floor((float64(a.vars.mood))/20.0))
}
