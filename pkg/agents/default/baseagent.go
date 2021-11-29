package baseagent

type BaseAgent struct {
	//don't direclty access or modify these during runtime pls
	HP    int
	Floor int
	// uuid
	//personality Personality
}

type Agent interface {
	Run(t *Tower) // maybe include tick Time (can also be sim phases)
}

// func (bA *BaseAgent) HPDecay() {
// 	//Defines amount of time a given agent can survive without food
// 	//agents could be assigned a random decay function at creation?
// 	//defined by infra
// 	bA.hunger -= N //this could be default behaviour
// }
// func (bA *BaseAgent) Communicate(floor type_with_3_vals, msg Message) {
// 	//if floor == UP
// 	//  broadcast msg to everyone on yur floor and the one above
// 	//if floor == DOWN
// 	//  broadcast msg to everyone on yur floor and the one below
// 	//if floor == SAME
// 	//  broadcast message to everyone on only your floor
// }

// func (bA *BaseAgent) ObserveTowerState(t Tower) { // <big ass list of stuff that the agent is allowed to see from the tower>{
// }

// func (bA *BaseAgent) MemoryDecay() //needs to be implemented later, assume agents can remeber forever atm? the agent teams will handle how to remember stuff
