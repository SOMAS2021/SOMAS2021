package team2

import (
	"testing"
	"math"
)

func TestCalcReward(t *testing.T) {
	overEating := 86 // this should be updated if param changed in reward.go
	type calcRewardTest struct {
		arg1, arg2, arg3, arg4, arg5, arg6, arg7 int
		expected1 float64
	}

	var calcRewardTests = []calcRewardTest{
		calcRewardTest{}
	}

	/*for _, test := range initTableTests{
		output := InitTable(test.arg1, test.arg2)
			t.Errorf("Length %d not equal to expected %d", len(output[0]), test.expected2)
	}*/
}
