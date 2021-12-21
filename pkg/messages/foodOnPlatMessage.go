package messages


type foodOnPlatMessage struct{
  *baseMessage
  food int
}

func NewfoodOnPlatMessage(SenderFloor int, food int) *foodOnPlatMessage{
  msg:= &foodOnPlatMessage{
    baseMessage: NewBaseMessage(SenderFloor),
    food:food,
  }
  return msg
}

func(msg foodOnPlatMessage) MessageType() string{
  return "foodOnPlatMessage"
}
