package team3

import (
	"math"
)

/*
stubbornness range = 0-100
mood range = 0-100
morality range = 0-100
*/

func takeFoodCalculation(a *CustomAgent3) int64 {
	//food taken is: 10 - moarlity/20 - mood/25
	var try_take_food = 10.0 - math.Floor(float64(a.vars.morality/20.0)) - math.Floor(float64(a.vars.mood)/20.0)
	return int64(try_take_food)
}
