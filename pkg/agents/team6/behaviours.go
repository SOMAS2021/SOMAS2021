package team6

import (
	"math"
)

func (a *CustomAgent6) updateBehaviour() {

	aConf := a.config
	behaviourMax, behaviourMin := a.getBehaviourRange()

	hpScore := 1 - float64(a.HP())/100.0 //map (minHP,maxHP) -> (1,0)
	floor := a.Floor() + 1
	// floor = 600 //use for debug

	if floor > a.maxFloorGuess {
		a.maxFloorGuess = floor + 1
	}

	fdNum := float64(floor) / float64(a.maxFloorGuess) * aConf.lambda
	fdDen := float64(aConf.lambda)
	floorScore := math.Exp(fdNum) / math.Exp(fdDen)

	behaviourParams := []float64{hpScore, floorScore}

	weights := aConf.paramWeights
	behaviourPrediction := 0.0

	for i := range weights {
		behaviourPrediction += behaviourParams[i] * weights[i]
	}

	updateDir := behaviour(behaviourPrediction)*aConf.maxBehaviourThreshold - a.currBehaviour //find direction we need to move in to reach new behaviour prediction
	updateMag := updateDir * behaviour(1-aConf.stubbornness)                                  //scale movement by stubbornness (minStubborn, maxStubborn) -> (fullMovement, 0)
	newBehaviour := a.currBehaviour + updateMag

	if newBehaviour > behaviourMax { //limit behaviour to max swing
		a.currBehaviour = behaviourMax
	} else if newBehaviour < behaviourMin {
		a.currBehaviour = behaviourMin
	} else {
		a.currBehaviour = newBehaviour
	}

}

func (a *CustomAgent6) getBehaviourRange() (behaviourMax, behaviourMin behaviour) {

	aConf := a.config

	maxT := a.config.maxBehaviourThreshold

	bMax := behaviour(math.Min(float64(maxT), float64(aConf.baseBehaviour)+aConf.maxBehaviourSwing))
	bMin := behaviour(math.Max(0, float64(aConf.baseBehaviour)-aConf.maxBehaviourSwing))

	return bMax, bMin
}
