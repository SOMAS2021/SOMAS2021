package messages


type askFoodOnPlatMessage struct{
  *baseMessage
  food int
}

func NewaskFoodOnPlatMessage(SenderFloor int, food int) *askFoodOnPlatMessage{
  msg:= &askFoodOnPlatMessage{
    baseMessage: NewBaseMessage(SenderFloor),
    food:food,
  }
  return msg
}

func(msg askFoodOnPlatMessage) MessageType() string{
  return "askFoodOnPlatMessage"
}
