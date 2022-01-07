package team4EvoAgent

import (
	"math"

	"github.com/SOMAS2021/SOMAS2021/pkg/infra"
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

//TODO:
//Propose our own new treaties
//gotta update trustscore using treaty stuff this will just be if someone breaks the treaty and number of signatures reduce
//actually implement treaty conditions for take food. gotta check for conditions around TakeFood function.

/*------------------------OTHER AGENTS TREATY PROPOSAL HANDLING------------------------*/

// CONDITION OF AVALIABLE FOOD
// If available food <= critical amount * 2 then we reject.
// Reject a treaty if the coniditonop of available food is < or  <= because we don't want to be in a position of loss in any case DONE
// Reject all treaties where avalaible - request leave food (both amount and percentage) <= critical amount *2.

// CONDITION OF HP
// Reject all treaties that have condition HP <= critcal * 2 cause we ourselves are desperate at that point DONE
// Reject all treaties that involve condition HP < or < = because we don't want to be in a position of loss in any case DONE

// CONDITION OF FLOOR
// Reject all treaties with the conditionop of FLOOR having > or >= because we don't want to be in a position of loss in any case DONE

// Reject all treaties with duration > max days in critical
// If request is leave food amount or leave percent food and we are asked to leave 100% of the food then reject cause we would need to eat.

func (a *CustomAgentEvo) HandleProposeTreaty(msg messages.ProposeTreatyMessage) {
	treaty := msg.Treaty()

	// Switch Case for rejecting treaties
	switch {
	case len(a.ActiveTreaties()) > 0: // If one treaty already signed then we reject any other incoming ones TODO: Implement more robust method of multi-treaty handling
		fallthrough
	case treaty.Request() == messages.Inform: //TODO no idea what inform actually does/how to use it
		fallthrough
	case treaty.Condition() == messages.HP && treaty.ConditionValue() <= a.HealthInfo().HPReqCToW*2:
		fallthrough
	case treaty.Condition() == messages.HP && (treaty.ConditionOp() == messages.LT || treaty.ConditionOp() == messages.LE): // Reject any case where they want us to accept treaties with less than any health %
		fallthrough
	case treaty.Request() == messages.LeavePercentFood && (treaty.RequestValue() >= 100 || treaty.RequestValue() <= 0):
		fallthrough
	case treaty.Condition() == messages.AvailableFood && treaty.ConditionValue() <= a.HealthInfo().HPReqCToW*2:
		fallthrough
	case treaty.Condition() == messages.AvailableFood && treaty.Request() == messages.LeaveAmountFood && treaty.ConditionValue()-treaty.RequestValue() <= a.HealthInfo().HPReqCToW*2:
		fallthrough
	case treaty.Condition() == messages.AvailableFood && treaty.Request() == messages.LeavePercentFood && int(float64(100-treaty.RequestValue()/100))*treaty.ConditionValue() <= a.HealthInfo().HPReqCToW*2:
		fallthrough
	case treaty.Condition() == messages.AvailableFood && (treaty.ConditionOp() == messages.LT || treaty.ConditionOp() == messages.LE):
		fallthrough
	case treaty.Condition() == messages.Floor && (treaty.ConditionOp() == messages.GT || treaty.ConditionOp() == messages.GE):
		fallthrough
	case a.HP() < a.HealthInfo().HPCritical: // Reject all treaties at that point in time when we're below critical health (all treaty requests are relevant to food only so this applies).
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

	// We propogate the treaty to the floor above/below
	a.propogateTreaty(msg)
}

/*------------------------REJECT A TREATY------------------------*/

func (a *CustomAgentEvo) rejectTreaty(msg messages.ProposeTreatyMessage) { // change later to not make sus
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
		a.UpdateFoodFromTreatyToAgent(treaty)
	}
}

func (a *CustomAgentEvo) handleActiveTreatyConditionsFloor(treaty messages.Treaty) {
	switch {
	case treaty.ConditionOp() == messages.EQ && a.Floor() == treaty.ConditionValue():
		fallthrough
	case treaty.ConditionOp() == messages.LT && a.Floor() < treaty.ConditionValue():
		fallthrough
	case treaty.ConditionOp() == messages.LE && a.Floor() <= treaty.ConditionValue():
		a.UpdateFoodFromTreatyToAgent(treaty)
	}
}

func (a *CustomAgentEvo) handleActiveTreatyConditionsAvailableFood(treaty messages.Treaty) {
	switch {
	case treaty.ConditionOp() == messages.EQ && a.CurrPlatFood() == food.FoodType(treaty.ConditionValue()):
		fallthrough
	case treaty.ConditionOp() == messages.GT && a.CurrPlatFood() > food.FoodType(treaty.ConditionValue()):
		fallthrough
	case treaty.ConditionOp() == messages.GE && a.CurrPlatFood() >= food.FoodType(treaty.ConditionValue()):
		a.UpdateFoodFromTreatyToAgent(treaty)
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

//messages.LeaveFoodAmt GT- check if a.currPlatFood - intended < request : foodtaken = a.currPlatFloor - request - 1
//GE - foodtaken = a.currPlatFloor - request
//LT - foodtaken = a.currPlatFood - request + 1
//LE - foodtaken = a.currPlatFood - request
//EQ - foodtaken = a.currPlatFood - request

//messages.LeavePercentFood GT - check if a.currPlatFood - intended < a.currentPlatFood * foodpercen/100 : foodtaken = a.currPlatFloor - a.currplatfood * request/100  - 1
//GE - foodtaken = a.currPlatFloor - request/100 * a.currentPlatFood
//LT - foodtaken = a.currPlatFood - request/100 * a.currentPlatFood + 1
//LE - foodtaken = a.currPlatFood - request/100 * a.currentPlatFood
//EQ - foodtaken = a.currPlatFood - request/100 * a.currentPlatFood

/*------------CHANGE THE AMOUNT OF FOOD WE TAKE DEPENDING ON THE ACTIVE TREATY CONDITIONS ------------*/
func (a *CustomAgentEvo) UpdateFoodFromTreatyToAgent(treaty messages.Treaty) {
	switch {
	case treaty.Request() == messages.LeaveAmountFood:
		if treaty.RequestOp() == messages.EQ {
			a.params.intendedFoodToTake = a.CurrPlatFood() - food.FoodType(treaty.RequestValue())
		} else if treaty.RequestOp() == messages.GT && a.CurrPlatFood()-a.params.intendedFoodToTake < food.FoodType(treaty.RequestValue()) {
			a.params.intendedFoodToTake = food.FoodType(math.Max(0, float64(a.CurrPlatFood()-food.FoodType(treaty.RequestValue())-1))) //TODO: make sure intended amount is not -ve. just put a max between 0 and that equation.

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
