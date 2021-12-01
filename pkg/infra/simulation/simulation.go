package simulation

import (
	"log"

	baseagent "github.com/SOMAS2021/SOMAS2021/pkg/agents/default"
	tower "github.com/SOMAS2021/SOMAS2021/pkg/infra/tower"
	"github.com/divan/goabm/abm"
)

type SimEnv struct {
	FoodOnPlatform float64
	AgentCount     int
	AgentHP        int
	Iterations     int
}

func New(foodOnPlat float64, agentCount int, agentHP int, iterations int) *SimEnv {

	s := &SimEnv{
		FoodOnPlatform: foodOnPlat,
		AgentCount:     agentCount,
		AgentHP:        agentHP,
		Iterations:     iterations,
	}
	// can do other inits here
	return s
}

func (sE *SimEnv) Simulate() {

	a := abm.New()
	tower := tower.New(sE.FoodOnPlatform, 1, sE.AgentCount)
	a.SetWorld(tower)

	for i := 0; i < sE.AgentCount; i++ {

		agent, err := baseagent.New(a, i, sE.AgentHP)
		if err != nil {
			log.Fatal(err)
		}
		a.AddAgent(agent)
		tower.SetAgent(i, agent)

	}

	a.LimitIterations(sE.Iterations)
	a.StartSimulation()

}

// Draft for WIP

/*

//more parameters to be considered
// What rules do the agents collectively agree on before they enter the tower.
// Do agents that replace dead agents know about these pre agreed rules?
// Number of messages that can be sent during a turn (time between the movement of the platform) or time taken for given message types to be sent
// Time taken for agents to respond to messages



func communication(agents *[]BaseAgent, messages *[]Message, commFloor int) {

	// allocate curr Messages

	// collect messages from agents

	// put new messages in inboxes

}

func eating(currAgent *BaseAgent, foodLeft *float64, maxHealthPoints int) {

	// first eating - ask agent to return food taken

	// second - updated health of curr Agent

}

func death(agents *[]BaseAgent) {

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
	// var minHealthPoints int = 0
	var numberOfAgents uint64
	var totalFood float64

	foodRange := foodToSatisfy - foodToSatisfice
	totalFood = foodToSatisfice + rand.Float64()*foodRange

	//-- instantiate agents here !! TALK ABOUT THIS WITH THE GROUP

	//-- HP decay function

	//tower params (- could be set in tower and called.)
	var numFloors uint64
	var agentsPerFloor uint64
	var currFloor int

	numFloors = numberOfAgents / agentsPerFloor

	var environment = simEnv{agentsPerFloor, foodLeft, ticksSpentByPlatformOnFloor}

	var agentListTemp []BaseAgent = []BaseAgent{{HP: 12, Floor: 2}}

	var tower = Tower{FoodOnPlatform: foodLeft, FloorOfPlatform: 0, Agents: agentListTemp}

	for currentDay < simLength {
		// ticks for loop TODO: import time package
		// time.tick()
		var tickTmp int = 0
		currFloor = 0
		foodLeft = totalFood
		currentDay++

		//message instantiation
		var messages []Message

		for tickTmp < ticksPerDay {

			//wake up (do we need to define agents sleeping?)

			// communicate phase
			communication(&tower.Agents, &messages, commFloors)
			// every tick/increment - communication function passes through the outbox
			// of every agent and checks whether there is a message 'waiting' to be sent
			// (this wait would be a very short period of time == one tick)
			// if ther a message to be sent - it is allocated to the right agent to be
			// received on the next tick/increment

			if currFloor != int(numFloors) {
				// platform move phase
				eating(&tower.Agents[int(currFloor)], &foodLeft, maxHealthPoints)
				//do we need to reset platform or just pass the top floor in the eating function
				currFloor++
			}

		}
		death(&tower.Agents)
		// when do we replace agents, as soon as they die or on the reshuffle?
		replace(&simEnv)
		daysUntilReshuffle--
		if daysUntilReshuffle == 0 {
			reshuffle(&tower.Agents, agentsPerFloor)
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

*/
