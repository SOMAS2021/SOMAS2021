package infra

import (
	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
	"github.com/google/uuid"
)

func (a *Base) ReceiveMessage() messages.Message {
	select {
	case msg := <-a.inbox:
		return msg
	default:
		return nil
	}
}

func (a *Base) SendMessage(msg messages.Message) {
	a.tower.SendMessage(a.floor, msg)
}

func (a *Base) ActiveTreaties() map[uuid.UUID]messages.Treaty {
	return a.activeTreaties
}

func (a *Base) AddTreaty(treaty messages.Treaty) {
	a.activeTreaties[treaty.ID()] = treaty
}

func (a *Base) DeleteTreaty(treatyID uuid.UUID) {
	delete(a.activeTreaties, treatyID)
}

func (a *Base) HandleAskHP(msg messages.AskHPMessage) {
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), a.HP())
	a.SendMessage(reply)
	a.Log("I recieved an askHP message from ", Fields{"senderFloor": msg.SenderFloor(), "myFloor": a.Floor()})

}
func (a *Base) HandleAskFoodTaken(msg messages.AskFoodTakenMessage)                      {}
func (a *Base) HandleAskIntendedFoodTaken(msg messages.AskIntendedFoodIntakeMessage)     {}
func (a *Base) HandleRequestLeaveFood(msg messages.RequestLeaveFoodMessage)              {}
func (a *Base) HandleRequestTakeFood(msg messages.RequestTakeFoodMessage)                {}
func (a *Base) HandleResponse(msg messages.BoolResponseMessage)                          {}
func (a *Base) HandleStateFoodTaken(msg messages.StateFoodTakenMessage)                  {}
func (a *Base) HandleStateHP(msg messages.StateHPMessage)                                {}
func (a *Base) HandleStateIntendedFoodTaken(msg messages.StateIntendedFoodIntakeMessage) {}

// Note: You can override this function depending on how you want to handle treaties.
// This implementation automatically accepts treaties. You probably don't want to do this.
func (a *Base) HandleProposeTreaty(msg messages.ProposeTreatyMessage) {
	treaty := msg.Treaty()
	treaty.SignTreaty()
	a.activeTreaties[msg.TreatyID()] = treaty
	reply := msg.Reply(a.ID(), a.Floor(), msg.SenderFloor(), true)
	a.SendMessage(reply)
	a.Log("Accepted treaty", Fields{"proposerID": msg.SenderID(), "proposerFloor": msg.SenderFloor(),
		"treatyID": msg.TreatyID()})
}

// Note: You can override this function depending on how you want to handle treaties.
// This implementation automatically increments the signature count of the treaty if it was accepted.
func (a *Base) HandleTreatyResponse(msg messages.TreatyResponseMessage) {
	if msg.Response() {
		treaty := a.activeTreaties[msg.TreatyID()]
		treaty.SignTreaty()
		a.activeTreaties[msg.TreatyID()] = treaty
	}
}

func (a *Base) HandlePropogate(msg messages.Message) {
	a.SendMessage(msg)
}
