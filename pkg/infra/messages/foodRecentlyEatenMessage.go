package messages

//The last time an agent eats this is how much food they ate
type FoodRecentlyEatenMessage struct {
	*baseMessage
	FoodAmt uint
}
