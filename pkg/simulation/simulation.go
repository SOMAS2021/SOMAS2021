//Ask people to be nice and to not break the simulation rules!!!!!!!!!
type simEnv struct {
	agents_per_floor           int
	platformStatingFood        int
	timeSpentByPlatformOnFloor int
}

//more parameters to be considered
// What rules do the agents collectively agree on before they enter the tower.
// Do agents that replace dead agents know about these pre agreed rules?
// Number of messages that can be sent during a turn (time between the movement of the platform) or time taken for given message types to be sent
// Time taken for agents to respond to messages

func (sE *simEnv) Replace( /*add relevant inputs*/ ) /*add relevant return values*/ {
	//implementation
}

func (eE *simEnv) reshuffle() {

}

// by making stuff lower case, they dont get access to it in other files, oh baby yeah lets do this
func (sE *simEnv) Simulate() {
	//come up with rules
	// Initialisation Phase - set intial parameters of sim
	// time parameters
	var simLength int
	var currentDay int = 0 // will not change throughout simulations
	var ticksPerDay int
	var commFloors int // +/- floors that can be communicated with
	var daysUntilReshuffle int
	var daysPerReshuffle int

	ticksPerDay = 500
	daysPerReshuffle = 10
	daysUntilReshuffle = daysPerReshuffle

	// platform parameters
	var foodLeft float64
	var ticksSpentByPlatformOnFloor int = 1

	//agent params
	var foodToSatisfice float64 // don't loose HP, dont gain HP
	var foodToSatisfy float64   // gain HP?
	var maxHealthPoints int = 100
	var minHealthPoints int = 0
	var numberOfAgents int
	var totalFood float64

	totalFood = foodToSatisfice * numberOfAgents
	// replace foodToSatisfice with rand(foodToSatisfice, foodToSatisfy)

	//-- instantiate agents here !! TALK ABOUT THIS WITH THE GROUP

	//-- HP decay function

	//tower params (- could be set in tower and called.)
	var numFloors int
	var agentsPerFloor int
	var currFloor int

	numFloors = numberOfAgents / agentsPerFloor

	var environment = simEnv{agentsPerFloor, foodLeft, ticksSpentByPlatformOnFloor}

	var tower = Tower{foodLeft, topFloor, agents}

	for currentDay; currentDay < simLength; currentDay++ {
		// ticks for loop TODO: import time package
		// time.tick()
		var tickTmp int = 0
		currFloor = numFloors
		foodLeft = totalFood

		for tickTmp < ticksPerDay {
			//wake up (do we need to define agents sleeping?)

			// communicate phase
			communication(tower.agents, messages, tickTmp)
			// every tick/increment - communication function passes through the outbox
			// of every agent and checks whether there is a message 'waiting' to be sent
			// (this wait would be a very short period of time == one tick)
			// if ther a message to be sent - it is allocated to the right agent to be
			// received on the next tick/increment

			// platform move phase
			eating(tower.agents, currFloor, agentsPerFloor, maxHealthPoints)
			//do we need to reset platform or just pass the top floor in the eating function
			currFloor--

		}
		death(tower.agents)
		// when do we replace agents, as soon as they die or on the reshuffle?
		replace(environment)
		daysUntilReshuffle--
		if daysUntilReshuffle == 0 {
			reshuffle(tower.agents, agentsPerFloor)
			daysUntilReshuffle = daysPerReshuffle
		}
		// print data requested by the frontend team
	}

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