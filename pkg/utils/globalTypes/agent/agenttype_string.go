// Code generated by "stringer -type=AgentType"; DO NOT EDIT.

package agent

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Team1Agent1-1]
	_ = x[Team1Agent2-2]
	_ = x[Team2-3]
	_ = x[Team3-4]
	_ = x[Team4-5]
	_ = x[Team5-6]
	_ = x[Team6-7]
	_ = x[Team7-8]
	_ = x[RandomAgent-9]
}

const _AgentType_name = "Team1Agent1Team1Agent2Team2Team3Team4Team5Team6Team7RandomAgent"

var _AgentType_index = [...]uint8{0, 11, 22, 27, 32, 37, 42, 47, 52, 63}

func (i AgentType) String() string {
	i -= 1
	if i < 0 || i >= AgentType(len(_AgentType_index)-1) {
		return "AgentType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _AgentType_name[_AgentType_index[i]:_AgentType_index[i+1]]
}