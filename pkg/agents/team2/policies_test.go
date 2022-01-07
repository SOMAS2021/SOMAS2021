package team2

import (
	"testing"
)

func TestInitPolicies(t *testing.T) {
	type initPoliciesTest struct {
		arg1, arg2, expected1, expected2 int
		expected3 float64
	}

	var initPoliciesTests = []initPoliciesTest{
		initPoliciesTest{3, 3, 3, 3, 1.0/3.0},
		initPoliciesTest{10, 1, 10, 1, 1.0/1.0},
		initPoliciesTest{50, 100, 50, 100, 1.0/100.0},
	}

	// Check dimensions of policy table are correct and check uniform distribution is correct
	for _, test := range initPoliciesTests{
		output := InitPolicies(test.arg1, test.arg2)
		if len(output) != test.expected1 {
			t.Errorf("Length %d not equal to expected %d", len(output), test.expected1)
		}
		for i := 0; i < test.arg1; i++ {
			if len(output[i]) != test.expected2 {
				t.Errorf("Length %d not equal to expected %d", len(output[i]), test.expected2)
			}
			for j := 0; j < test.arg2; j++ {
				if output[i][j] != test.expected3 {
					t.Errorf("Uniform probability %f not equal to expected %f", output[i][j], test.expected3)
				}
			}
		}
	}

}

func TestUpdatePolicies(t *testing.T) {
	// TODO: Implement this test if considered appropriate
}

func TestAdjustPolicies() {
	// TODO: Implement this test if considered appropriate
}
