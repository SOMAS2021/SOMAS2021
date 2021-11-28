type Agent interface{
    Run(t *Tower) // maybe include tick Time (can also be sim phases)
}
type baseAgent struct {
    //don't direclty access or modify these during runtime pls
    HP 
    uuid
    floor
    //personality Personality
}
type seflishAgent struct{
    bA baseAgent 
    const selfishness float64 = 1.0
    personality Personality
}
func (sA *selfishAgent)Run(){
    amt  = calculateSatisfaction_amt()
    tower.TakeFood(amt)
    Communicate(1, "hi")
}
type selflessAgent struct{
    bA baseAgent 
    const selfishness float64 = 0.0
    personality Personality
}
type randomAgent struct{
    bA baseAgent 
    selfishness float64 //randomly assigned at construction
    personality Personality
}
func (bA *baseAgent) HPDecay(){
    //Defines amount of time a given agent can survive without food
    //agents could be assigned a random decay function at creation?
    //defined by infra
    bA.hunger -= N //this could be default behaviour
}
func (bA *baseAgent) Communicate(floor type_with_3_vals, msg Message){
    //if floor == UP
    //  broadcast msg to everyone on yur floor and the one above
    //if floor == DOWN
    //  broadcast msg to everyone on yur floor and the one below
    //if floor == SAME
    //  broadcast message to everyone on only your floor
}

func (bA *baseAgent) ObserveTowerState(t Tower) <big ass list of stuff that the agent is allowed to see from the tower>{
} 

func (bA *baseAgent) MemoryDecay() //needs to be implemented later, assume agents can remeber forever atm? the agent teams will handle how to remember stuff
