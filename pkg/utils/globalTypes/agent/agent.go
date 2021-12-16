package agent

type Agent interface {
	Run()
	IsAlive() bool
}
