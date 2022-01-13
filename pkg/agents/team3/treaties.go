package team3

import (
	"math/rand"

	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/google/uuid"
)

//AgentPosition signifies how strong it needs to be to approve a treaty.
type AgentPosition int

//FoodTaken signifies how much food we need to take to approve a treaty.
type FoodTaken int

const (
	Strong AgentPosition = iota + 1
	Healthy
	Average
	Weak
	SurvivalLevel
	Reject
)

const (
	VeryLarge FoodTaken = iota + 1
	Large
	Moderate
	Little
	SurvivalAmount
	TooLittle
)

func (a *CustomAgent3) treatyFull() bool {
	return len(a.ActiveTreaties()) > 0
}

func (a *CustomAgent3) treatyPendingResponse() bool {
	return !(a.knowledge.treatyProposed.ProposerID() == uuid.Nil)
}

func (a *CustomAgent3) proposeTreatiesImmoral(floor int) { //troll treaties (then add some more b)
	tr := messages.NewTreaty(messages.Floor, a.Floor()-1, messages.LeavePercentFood, 99, messages.GT, messages.GT, a.knowledge.reshuffleEst/2, a.ID())
	a.knowledge.treatyProposed = *tr //remember the treaty we proposed
	msg := messages.NewProposalMessage(a.BaseAgent().ID(), a.Floor(), floor, *tr)
	a.SendMessage(msg)
	a.Log("I sent a treaty")
}

func (a *CustomAgent3) proposeTreatiesMoral(direction int) {
	tr := messages.NewTreaty(messages.HP, 20, messages.LeavePercentFood, 95, messages.GT, messages.GT, 3, a.ID())
	r := rand.Intn(3)
	switch r {
	case 0:
		tr = messages.NewTreaty(messages.HP, 60, messages.LeavePercentFood, 60, messages.GT, messages.GT, 5, a.ID())
	case 1:
		tr = messages.NewTreaty(messages.HP, 10, messages.LeavePercentFood, 95, messages.GT, messages.GT, 10, a.ID())
	case 2:
		tr = messages.NewTreaty(messages.AvailableFood, 50, messages.LeaveAmountFood, a.foodReqCalc(a.HP(), a.HealthInfo().HPReqCToW), messages.LT, messages.GT, 5, a.ID())
	}
	a.knowledge.treatyProposed = *tr //remember the treaty we proposed
	msg := messages.NewProposalMessage(a.BaseAgent().ID(), a.Floor(), a.Floor()-direction, *tr)
	a.SendMessage(msg)
	a.Log("I sent a treaty")
}

// Returns the AgentPosition (relative strength measure) of the agent when at the minimum HP defined by the condition and conditionOp
func (a *CustomAgent3) requiredHPLevel(treaty messages.Treaty) AgentPosition {
	if treaty.ConditionOp() == messages.LT || treaty.ConditionOp() == messages.LE || treaty.ConditionValue() > 100 {
		return SurvivalLevel
	}
	switch hp := treaty.ConditionValue(); {
	case hp >= 75:
		return Strong
	case hp >= 55:
		return Healthy
	case hp >= 35:
		return Average
	case hp >= a.HealthInfo().WeakLevel:
		return Weak
	case hp == a.HealthInfo().HPCritical:
		return SurvivalLevel
	default:
		return Reject
	}
}

// Determining if a given floor means an agent is in a good / bad position relies on knowledge that the agent has no access to.
// Hence, initial approach is to assume the treaty is risky, and could possibly activate when the agent is at SurvivalLevel.
func (a *CustomAgent3) requiredFloorLevel(treaty messages.Treaty) AgentPosition {
	return Reject
}

// Same as Floor, initial approach is to assume the treaty is risky, and could possibly activate when the agent is at SurvivalLevel.
func (a *CustomAgent3) requiredAvailFoodLevel(treaty messages.Treaty) AgentPosition {
	return Reject
}

// Calculates food available to eat if request applied to current platform food, and uses this as an estimate for the general case.
func (a *CustomAgent3) reqFoodTakenEstimate(treaty messages.Treaty, percentage bool) FoodTaken {
	foodToEatCalc := 0

	if treaty.RequestOp() == messages.LT || treaty.RequestOp() == messages.LE || treaty.RequestOp() == messages.EQ {
		return VeryLarge
	}

	if percentage {
		foodToEatCalc = int(float64(a.CurrPlatFood()) * float64((100.0-float64(treaty.RequestValue()))/100.0))
	} else {
		foodToEatCalc = int(int(a.CurrPlatFood()) - treaty.RequestValue())
	}
	switch foodToEat := foodToEatCalc; {
	case foodToEat > a.foodReqCalc(85, 85):
		return VeryLarge
	case foodToEat > a.foodReqCalc(60, 60):
		return Large
	case foodToEat > a.foodReqCalc(40, 40):
		return Moderate
	case foodToEat > a.foodReqCalc(a.HealthInfo().WeakLevel+15, a.HealthInfo().WeakLevel+15):
		return Little
	case foodToEat > a.foodReqCalc(a.HealthInfo().WeakLevel, a.HealthInfo().WeakLevel):
		return SurvivalAmount
	default:
		return TooLittle
	}
}

//HandleProposeTreaty 1. requiredAgentPosition evaluates the condition, 2. foodTakenEstimate evaluates the request,
// 3. agentVarsPassed uses agent params with evaluations, 4. Reply sent which accepts/rejects the treaty
func (a *CustomAgent3) HandleProposeTreaty(msg messages.ProposeTreatyMessage) {
	if a.treatyFull() {
		reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), false)
		a.SendMessage(reply)
	} else {

		treaty := msg.Treaty()
		minActivationLevel := SurvivalLevel
		response := false

		switch treaty.Condition() {
		case messages.HP:
			minActivationLevel = a.requiredHPLevel(treaty)
		case messages.Floor:
			minActivationLevel = a.requiredFloorLevel(treaty)
		case messages.AvailableFood:
			minActivationLevel = a.requiredAvailFoodLevel(treaty)
		}

		foodTakenEstimate := a.reqFoodTakenEstimate(treaty, treaty.Request() == messages.LeavePercentFood)

		// If agent is in a bad mood, it will only accept treaties that take effect when it is in a strong position.
		// If agent has low morality, it will only accept treaties that involve it taking large amounts of food.
		agentVarsPassed := a.vars.mood > (20*int(minActivationLevel)-20) && a.vars.morality > (20*int(foodTakenEstimate)-20) && a.vars.morality < (20*int(foodTakenEstimate)+20)
		// Check duration is not too long, and use stubbornness to decide if an agent gives in and accepts treaties with at least 5 signatures
		allChecksPassed := agentVarsPassed && treaty.Duration() < 2*a.knowledge.reshuffleEst
		if treaty.SignatureCount() >= 5 && !allChecksPassed {
			allChecksPassed = a.read()
		}

		//use agent variables, foodTakenEstimate, and requiredAgentPosition to accept/reject
		if allChecksPassed && treaty.Request() != messages.Inform { // dont accept HP inform requests
			response = true
			treaty.SignTreaty()
			a.AddTreaty(treaty)
		} else { // reject other treaties
			response = false
		}
		reply := msg.Reply(a.BaseAgent().ID(), a.Floor(), msg.SenderFloor(), response)
		a.SendMessage(reply)
	}
}

//HandleTreatyResponse handles the response of a treaty
func (a *CustomAgent3) HandleTreatyResponse(msg messages.TreatyResponseMessage) {
	if msg.Response() {
		if a.knowledge.treatyProposed.ID() != uuid.Nil { //check in case something went wrong.
			treaty := a.knowledge.treatyProposed //signed our proposed treatry
			treaty.SignTreaty()
			a.AddTreaty(treaty)
		}
	}
	a.knowledge.treatyProposed = *messages.NewTreaty(messages.HP, 0, messages.LeaveAmountFood, 0, messages.GT, messages.GT, 0, uuid.Nil) //restart the sent treaty.
}
