package Tower

import (
	. "github.com/SOMAS2021/SOMAS2021/pkg/agents/default"
)

type Tower struct {
	FoodOnPlatform  float64
	FloorOfPlatform uint64
	Agents          []BaseAgent
}

// func (t Tower) TakeFood(a Agent, amt uint64) {
// 	if foodPlatform_floor == a.floor {
// 		if amt > tower.PlatformFood {
// 			a.hunger += PlatformFood
// 			t.PlatformFood = 0
// 			fmt.Println("Tried to take more food than exists, were given all the food")
// 		} else {
// 			a.hunger += amt
// 			t.PlatformFood -= amt
// 		}
// 	} else {
// 		fmt.Println("Platform not on your floor, no food for you")
// 	}
// }
