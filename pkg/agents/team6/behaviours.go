package team6

import (
	"math"
)

// Updates agent social motive:
// First updates behaviour weights, then computes behaviour change based on weights and input parameters (HP, floor)
func (a *CustomAgent6) updateBehaviour() {
	a.updateBehaviourWeights()
	aConf := a.config
	behaviourMax, behaviourMin := a.behaviourRange()

	// Utility func for HP
	hpScore := 1 - float64(a.HP())/float64(a.HealthInfo().MaxHP)
	floor := a.Floor() + 1

	// Update agent's guess of how many floors there are in the tower
	if floor > a.maxFloorGuess {
		a.maxFloorGuess = floor + 1
	}

	// Utility func for floor
	fdNum := float64(floor) / float64(a.maxFloorGuess) * aConf.lambda
	fdDen := float64(aConf.lambda)
	floorScore := math.Exp(fdNum) / math.Exp(fdDen)

	// Compuute predicted new behaviour (with no constraints) by taking dot product between weights vector and utility funcs vector
	// Normalised between 0 and 1
	weights := aConf.paramWeights
	behaviourPrediction := hpScore*weights.HPWeight + floorScore*weights.floorWeight

	// Find new direction required to reach new behaviour prediction
	// Unconstrained new behaviour - current behaviour
	updateDir := behaviour(behaviourPrediction)*aConf.maxBehaviourThreshold - a.currBehaviour

	// Apply constraints
	// Scale movement by stubbornness (minStubborn, maxStubborn) -> (fullMovement, 0)
	updateMag := updateDir * behaviour(1-aConf.stubbornness)
	newBehaviour := a.currBehaviour + updateMag

	// Clip new behaviour between allowable behaviour range (based on behaviour swing)
	if newBehaviour > behaviourMax { //limit behaviour to max swing
		a.currBehaviour = behaviourMax
	} else if newBehaviour < behaviourMin {
		a.currBehaviour = behaviourMin
	} else {
		a.currBehaviour = newBehaviour
	}
}

// Returns range of allowable behaviours based on base behaviour and max behaviour swing
func (a *CustomAgent6) behaviourRange() (behaviourMax, behaviourMin behaviour) {
	aConf := a.config

	maxT := aConf.maxBehaviourThreshold

	bMax := behaviour(math.Min(float64(maxT), float64(aConf.baseBehaviour)+aConf.maxBehaviourSwing))
	bMin := behaviour(math.Max(0, float64(aConf.baseBehaviour)-aConf.maxBehaviourSwing))

	return bMax, bMin
}

// Updates weights on HP and floor based on current weight and situation
// If low on health and HP not weighted highly, then increase HP weight and decrease floor weight
// If low on food intake and floor not weighted highly, then increase floor weight and decrease HP weight
func (a *CustomAgent6) updateBehaviourWeights() {
	weights := &a.config.paramWeights
	if a.HP() < 20 && weights.HPWeight <= 0.95 {
		weights.HPWeight += 0.05
		weights.floorWeight -= 0.05
	}
	if a.averageFoodIntake < 1 && weights.floorWeight <= 0.9 {
		weights.HPWeight -= 0.1
		weights.floorWeight += 0.1
	}
}
