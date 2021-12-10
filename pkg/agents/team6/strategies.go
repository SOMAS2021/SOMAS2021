package team6

type thresholdData struct {
	satisficeThresh float64
	satisfyThresh   float64
	maxIntake       float64
}

func (a *CustomAgent6) foodIntake() float64 {
	thresholds := thresholdData{satisficeThresh: 20.0, satisfyThresh: 60.0, maxIntake: 80.0}

	switch a.currBehaviour.String() {
	case "Altruist":
		return 0.0
	case "Collectivist":
		return thresholds.satisficeThresh
	case "Selfish":
		return thresholds.satisfyThresh
	case "Narcissist":
		return thresholds.maxIntake
	default:
		return 0.0
	}
}
