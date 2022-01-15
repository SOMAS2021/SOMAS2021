package agent

type AgentType int

//go:generate go run golang.org/x/tools/cmd/stringer -type=AgentType
const (
	Team1 AgentType = iota + 1
	Team2
	Team3
	Team4
	Team5
	Team6
	Team7
	RandomAgent
)
