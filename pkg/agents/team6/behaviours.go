package team6

import (
	"math"
)

func (a *CustomAgent6) updateBehaviour() {
	a.updateBehaviourWeights()
	aConf := a.config
	behaviourMax, behaviourMin := a.behaviourRange()

	hpScore := 1 - float64(a.HP())/100.0
	floor := a.Floor() + 1

	if floor > a.maxFloorGuess {
		a.maxFloorGuess = floor + 1
	}

	fdNum := float64(floor) / float64(a.maxFloorGuess) * aConf.lambda
	fdDen := float64(aConf.lambda)
	floorScore := math.Exp(fdNum) / math.Exp(fdDen)

	weights := aConf.paramWeights
	behaviourPrediction := hpScore*weights.HPWeight + floorScore*weights.floorWeight

	// Find new direction required to reach new behaviour prediction
	updateDir := behaviour(behaviourPrediction)*aConf.maxBehaviourThreshold - a.currBehaviour
	// Scale movement by stubbornness (minStubborn, maxStubborn) -> (fullMovement, 0)
	updateMag := updateDir * behaviour(1-aConf.stubbornness)
	newBehaviour := a.currBehaviour + updateMag

	if newBehaviour > behaviourMax { //limit behaviour to max swing
		a.currBehaviour = behaviourMax
	} else if newBehaviour < behaviourMin {
		a.currBehaviour = behaviourMin
	} else {
		a.currBehaviour = newBehaviour
	}

}

func (a *CustomAgent6) behaviourRange() (behaviourMax, behaviourMin behaviour) {
	aConf := a.config

	maxT := aConf.maxBehaviourThreshold

	bMax := behaviour(math.Min(float64(maxT), float64(aConf.baseBehaviour)+aConf.maxBehaviourSwing))
	bMin := behaviour(math.Max(0, float64(aConf.baseBehaviour)-aConf.maxBehaviourSwing))

	return bMax, bMin
}

func (a *CustomAgent6) updateBehaviourWeights() {
	weights := &a.config.paramWeights
	if a.HP() < 20 && weights.HPWeight <= 0.95 {
		weights.HPWeight += 0.05
		weights.floorWeight -= 0.05
	}
	if a.averageFoodIntake < 10 && weights.floorWeight <= 0.9 {
		weights.HPWeight -= 0.1
		weights.floorWeight += 0.1
	}
}
