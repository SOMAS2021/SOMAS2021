package team4EvoAgent

import (
	"math"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

/*------------------------OTHER AGENTS TREATY PROPOSAL HANDLING------------------------*/

// CONDITION OF AVALIABLE FOOD
// If available food <= critical amount * 2 then we reject.
// Reject a treaty if the conditionop of available food is < or  <= because we don't want to be in a position of loss in any case
// Reject all treaties where avalaible - request leave food (both amount and percentage) <= critical amount *2.

// CONDITION OF HP
// Reject all treaties that have condition HP <= critcal * 2 cause we ourselves are desperate at that point
// Reject all treaties that involve condition HP < or < = because we don't want to be in a position of loss in any case

// CONDITION OF FLOOR
// Reject all treaties with the conditionop of FLOOR having > or >= because we don't want to be in a position of loss in any case

// Reject all treaties with duration > max days in critical
// If request is leave food amount or leave percent food and we are asked to leave 100% of the food then reject cause we would need to eat.

func (a *CustomAgentEvo) HandleProposeTreaty(msg messages.ProposeTreatyMessage) {
	treaty := msg.Treaty()

	// Switch Case for rejecting treaties
	switch {
	case len(a.ActiveTreaties()) > 0: // If one treaty already signed then we reject any other incoming ones
		fallthrough
	case treaty.Request() == messages.Inform: // No idea what inform actually does/how to use it
		fallthrough
	case treaty.Condition() == messages.HP && treaty.ConditionValue() <= a.HealthInfo().HPReqCToW*2:
		fallthrough
	case treaty.Condition() == messages.HP && (treaty.ConditionOp() == messages.LT || treaty.ConditionOp() == messages.LE): // Reject any case where they want us to accept treaties with less than any health %
		fallthrough
	case treaty.Request() == messages.LeavePercentFood && (treaty.RequestValue() >= 100 || treaty.RequestValue() <= 0):
		fallthrough
	case treaty.Condition() == messages.AvailableFood && treaty.ConditionValue() <= a.HealthInfo().HPReqCToW*2:
		fallthrough
	case treaty.Condition() == messages.AvailableFood && treaty.Request() == messages.LeaveAmountFood && treaty.ConditionValue()-treaty.RequestValue() <= a.HealthInfo().HPReqCToW*2: // want to leave at least buffer amount on plat
		fallthrough
	case treaty.Condition() == messages.AvailableFood && treaty.Request() == messages.LeavePercentFood && int(float64(100-treaty.RequestValue()/100))*treaty.ConditionValue() <= a.HealthInfo().HPReqCToW*2:
		fallthrough
	case treaty.Condition() == messages.AvailableFood && (treaty.ConditionOp() == messages.LT || treaty.ConditionOp() == messages.LE):
		fallthrough
	case treaty.Condition() == messages.Floor && (treaty.ConditionOp() == messages.GT || treaty.ConditionOp() == messages.GE):
		fallthrough
	case a.HP() < a.HealthInfo().HPCritical: // Reject all treaties at that point in time when we're below critical health (all treaty requests are relevant to food only so this applies).
		fallthrough
	case treaty.Condition() == messages.AvailableFood && treaty.Request() == messages.LeaveAmountFood && treaty.ConditionValue()-a.params.foodToEat["selfless"][2] <= treaty.RequestValue(): //the worst we want to be is in selfless
		fallthrough
	case treaty.Condition() == messages.AvailableFood && treaty.Request() == messages.LeavePercentFood && treaty.ConditionValue()-a.params.foodToEat["selfless"][2] <= treaty.ConditionValue()*treaty.RequestValue()/100: //the worst we want to be is in selfless
		fallthrough
	case treaty.Duration() >= a.HealthInfo().MaxDayCritical:
		fallthrough
	case int(a.params.globalTrust) < a.params.globalTrustLimits[0]:
		a.rejectTreaty(msg)
		return
	}

	// If it passes our conditions above then we accept the treaty
	treaty.SignTreaty()
	a.AddTreaty(treaty)
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), true)
	a.SendMessage(reply)
	a.Log("Accepted treaty", infra.Fields{"proposerID": msg.SenderID(), "proposerFloor": msg.SenderFloor(), "treatyID": msg.TreatyID()})
	// We propogate the treaty to the floor above/below if we have signed the treaty
	a.propogateTreaty(msg)
}

/*------------------------REJECT A TREATY------------------------*/

func (a *CustomAgentEvo) rejectTreaty(msg messages.ProposeTreatyMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), false)
	a.SendMessage(reply)
	a.Log("Rejected treaty", infra.Fields{"proposerID": msg.SenderID(), "proposerFloor": msg.SenderFloor(), "treatyID": msg.TreatyID()})
}

/*------------------------PROPOGATE A TREATY TO THE FLOOR ABOVE AND BELOW------------------------*/

func (a *CustomAgentEvo) propogateTreaty(msg messages.ProposeTreatyMessage) {
	treaty := msg.Treaty()
	toFloor := a.Floor() + 1
	if msg.SenderFloor() > a.Floor() {
		toFloor = a.Floor() - 1
	}
	propogateTreaty := messages.NewProposalMessage(a.ID(), a.Floor(), toFloor, treaty)
	a.SendMessage(propogateTreaty)
}

/*------------------------HANDLING OTHER AGENTS REPSONSES TO OUR TREATY PROPOSALS------------------------*/

func (a *CustomAgentEvo) HandleTreatyResponse(msg messages.TreatyResponseMessage) {
	if msg.Response() { // Check if there is a response
		_, ok := a.ActiveTreaties()[msg.TreatyID()] // Check their treaty is valid in our memory
		if ok {
			treaty := a.ActiveTreaties()[msg.TreatyID()] // Get the treaty from memory
			treaty.SetCount(treaty.SignatureCount() + 1) // Update the signature count to our accepted respondent
			a.ActiveTreaties()[msg.TreatyID()] = treaty  // Add it back into our activeTreaties
			a.addToGlobalTrust(a.params.trustCoefficients[0])
		}
	}
}

/*------------------------HANDLING OUR ACTIVE TREATIES ------------------------*/

func (a *CustomAgentEvo) handleActiveTreatyConditionsHP(treaty messages.Treaty) {
	switch {
	case treaty.ConditionOp() == messages.EQ && a.HP() == treaty.ConditionValue():
		fallthrough
	case treaty.ConditionOp() == messages.GT && a.HP() > treaty.ConditionValue():
		fallthrough
	case treaty.ConditionOp() == messages.GE && a.HP() >= treaty.ConditionValue():
		a.updateFoodFromTreatyToAgent(treaty)
	}
}

func (a *CustomAgentEvo) handleActiveTreatyConditionsFloor(treaty messages.Treaty) {
	switch {
	case treaty.ConditionOp() == messages.EQ && a.Floor() == treaty.ConditionValue():
		fallthrough
	case treaty.ConditionOp() == messages.LT && a.Floor() < treaty.ConditionValue():
		fallthrough
	case treaty.ConditionOp() == messages.LE && a.Floor() <= treaty.ConditionValue():
		a.updateFoodFromTreatyToAgent(treaty)
	}
}

func (a *CustomAgentEvo) handleActiveTreatyConditionsAvailableFood(treaty messages.Treaty) {
	switch {
	case treaty.ConditionOp() == messages.EQ && a.CurrPlatFood() == food.FoodType(treaty.ConditionValue()):
		fallthrough
	case treaty.ConditionOp() == messages.GT && a.CurrPlatFood() > food.FoodType(treaty.ConditionValue()):
		fallthrough
	case treaty.ConditionOp() == messages.GE && a.CurrPlatFood() >= food.FoodType(treaty.ConditionValue()):
		a.updateFoodFromTreatyToAgent(treaty)
	}
}

func (a *CustomAgentEvo) handleActiveTreatyConditions() {
	for _, treaty := range a.ActiveTreaties() {
		if treaty.SignatureCount() > 1 {
			if treaty.Condition() == messages.HP {
				a.handleActiveTreatyConditionsHP(treaty)
			} else if treaty.Condition() == messages.Floor {
				a.handleActiveTreatyConditionsFloor(treaty)
			} else if treaty.Condition() == messages.AvailableFood {
				a.handleActiveTreatyConditionsAvailableFood(treaty)
			}
		}
	}
}

/*------------CHANGE THE AMOUNT OF FOOD WE TAKE DEPENDING ON THE ACTIVE TREATY CONDITIONS ------------*/

//--------UpdateFoodFromTreatyToAgent explanation--------//

//messages.LeaveFoodAmt GT- check if a.currPlatFood - intended < request : foodtaken = a.currPlatFood - request - 1
//GE - foodtaken = a.currPlatFood - request
//LT - foodtaken = a.currPlatFood - request + 1
//LE - foodtaken = a.currPlatFood - request
//EQ - foodtaken = a.currPlatFood - request

//messages.LeavePercentFood GT - check if a.currPlatFood - intended < a.currentPlatFood * foodpercent/100 : foodtaken = a.currPlatFood - a.currplatfood * request/100  - 1
//GE - foodtaken = a.currPlatFood - (request/100 * a.currentPlatFood)
//LT - foodtaken = a.currPlatFood - (request/100 * a.currentPlatFood) + 1
//LE - foodtaken = a.currPlatFood - (request/100 * a.currentPlatFood)
//EQ - foodtaken = a.currPlatFood - (request/100 * a.currentPlatFood)

func (a *CustomAgentEvo) updateFoodFromTreatyToAgent(treaty messages.Treaty) {
	switch {
	case treaty.Request() == messages.LeaveAmountFood:
		if treaty.RequestOp() == messages.EQ {
			a.params.intendedFoodToTake = a.CurrPlatFood() - food.FoodType(treaty.RequestValue())
		} else if treaty.RequestOp() == messages.GT && a.CurrPlatFood()-a.params.intendedFoodToTake < food.FoodType(treaty.RequestValue()) {
			a.params.intendedFoodToTake = food.FoodType(math.Max(0, float64(a.CurrPlatFood()-food.FoodType(treaty.RequestValue())-1)))

		} else if treaty.RequestOp() == messages.GE && a.CurrPlatFood()-a.params.intendedFoodToTake <= food.FoodType(treaty.RequestValue()) {
			a.params.intendedFoodToTake = food.FoodType(math.Max(0, float64(a.CurrPlatFood()-food.FoodType(treaty.RequestValue()))))

		} else if treaty.RequestOp() == messages.LT && a.CurrPlatFood()-a.params.intendedFoodToTake > food.FoodType(treaty.RequestValue()) {
			a.params.intendedFoodToTake = food.FoodType(math.Max(0, float64(a.CurrPlatFood()-food.FoodType(treaty.RequestValue())+1)))

		} else if treaty.RequestOp() == messages.LE && a.CurrPlatFood()-a.params.intendedFoodToTake >= food.FoodType(treaty.RequestValue()) {
			a.params.intendedFoodToTake = food.FoodType(math.Max(0, float64(a.CurrPlatFood()-food.FoodType(treaty.RequestValue()))))
		}

	case treaty.Request() == messages.LeavePercentFood:
		if treaty.RequestOp() == messages.EQ {
			a.params.intendedFoodToTake = food.FoodType((int(a.CurrPlatFood()) * (1 - treaty.RequestValue()/100)))

		} else if treaty.RequestOp() == messages.GT && a.CurrPlatFood()-a.params.intendedFoodToTake < a.CurrPlatFood()*food.FoodType(treaty.RequestValue())/100 {
			a.params.intendedFoodToTake = food.FoodType((int(a.CurrPlatFood())*(1-treaty.RequestValue()/100) - 1))

		} else if treaty.RequestOp() == messages.GE && a.CurrPlatFood()-a.params.intendedFoodToTake <= a.CurrPlatFood()*food.FoodType(treaty.RequestValue())/100 {
			a.params.intendedFoodToTake = food.FoodType((int(a.CurrPlatFood()) * (1 - treaty.RequestValue()/100)))

		} else if treaty.RequestOp() == messages.LT && a.CurrPlatFood()-a.params.intendedFoodToTake > a.CurrPlatFood()*food.FoodType(treaty.RequestValue())/100 {
			a.params.intendedFoodToTake = food.FoodType((int(a.CurrPlatFood())*(1-treaty.RequestValue()/100) + 1))

		} else if treaty.RequestOp() == messages.LE && a.CurrPlatFood()-a.params.intendedFoodToTake >= a.CurrPlatFood()*food.FoodType(treaty.RequestValue())/100 {
			a.params.intendedFoodToTake = food.FoodType((int(a.CurrPlatFood()) * (1 - treaty.RequestValue()/100)))
		}
	}
}
