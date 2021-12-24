package agent

import "github.com/SOMAS2021/SOMAS2021/pkg/infra"

type Agent interface {
	Run()
	BaseAgent() *infra.Base
	IsAlive() bool
}
