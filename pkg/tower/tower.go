type Tower struct{
    FoodOnPlatform uint64
    FloorOfPlatform uint64
    agents []Agent
}

func (t Tower) TakeFood(a Agent, amt uint64){
    if foodPlatform_floor == a.floor {
        if amt > t.PlatformFood {
            a.hunger += t.PlatformFood
            t.PlatformFood = 0
            fmt.Println("Tried to take more food than exists, were given all the food")
        } else {
            a.hunger += amt
            t.PlatformFood -= amt
        }
    } else {
        fmt.Println("Platform not on your floor, no food for you")
    }
}
