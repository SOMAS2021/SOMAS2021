package messages

//The next time an agent eats this is how much food they will eat
type FoodToBeEatenMessage struct {
	*baseMessage
	FoodAmt uint
}