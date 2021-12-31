package infra

type FloorError struct{}
type NegFoodError struct{}
type AlreadyEatenError struct{}

func (e *FloorError) Error() string {
	return "platform is not on your floor"
}

func (e *NegFoodError) Error() string {
	return "cannot take negative food"
}

func (e *AlreadyEatenError) Error() string {
	return "already eaten today"
}
