//Ask people to be nice and to not break the simulation rules!!!!!!!!!
type simEnv struct{
    agents_per_floor int
    platformStatingFood int
    timeSpentByPlatformOnFloor int
}
//more parameters to be considered
// What rules do the agents collectively agree on before they enter the tower. 
// Do agents that replace dead agents know about these pre agreed rules? 
// Number of messages that can be sent during a turn (time between the movement of the platform) or time taken for given message types to be sent 
// Time taken for agents to respond to messages 

func (sE *simEnv) Replace(/*add relevant inputs*/) /*add relevant return values*/{
	//implementation
}
by making stuff lower case, they dont get access to it in other files, oh baby yeah lets do this
func (sE *simEnv) Simulate() {
    //come up with rules
    //Init Phase
    //set intial parameters of sim
    //-- sim_length (number of day)
    //-- agents communicating with agents K floors above and below = 1
    //-- ticks for the communicaiton phase
    //-- amount of food on platform
    
    //agent params
    //-- amount of food to satisfy
    //-- amount of food to satisfice
    //-- amount of HP that causes agents to die
    //-- HP decay function
    //tower params
    //-- number of floors
    //-- agents per floor = 1
    
    //- Instantiate simEnvironent with the parameters
    //- instantiate tower
    //- instantiate agents in the tower
    //day loop - number of days is defined as a parameter
        //wake up
        //loop until platform gets to bottom
            //communicate phase - assign number of ticks to this phase (run is executed every tick, (within a tick agents can do as much as they want)
            //move platform phase
        //reset platform
        //kill agents (check hunger, kill if hunger too high)
        //replace - parameter (when they die, when they shuffle)
        //random shuffle on Nth day - parameter
    
    //end of sim
    //- write to json file (could just dump the entire program state)
}