package main

import (
	sim "github.com/SOMAS2021/SOMAS2021/pkg/infra/simulation"
)

func main() {
	// can have frontend parameters come go straight into simEnv
	simEnv := sim.New(100, 3, 100, 10)
	simEnv.Simulate()
}
