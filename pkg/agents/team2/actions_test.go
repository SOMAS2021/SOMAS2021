package team2

import (
	"testing"
	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

func TestInitActionSpace(t *testing.T) {
	type initActionSpaceTest struct {
		arg1, expected int
	}

	var initActionSpaceTests = []initActionSpaceTest{
		initActionSpaceTest{0, 0},
		initActionSpaceTest{3, 3},
		initActionSpaceTest{99999, 99999},
	}
	for _, test := range initActionSpaceTests{
		// Check number of actions is correct
		if output := InitActionSpace(test.arg1); len(output.actionId) == test.expected {
		// Check that action IDs are correct
			for i := 0; i < test.arg1; i++ {
				if output.actionId[i] != i {
					t.Errorf("Output %d not equal to expected %d", len(output.actionId), test.expected)
				}
			}
		} else {
			t.Errorf("Output %d not equal to expected %d", len(output.actionId), test.expected)
		}
	}
}

func TestDisFood(t *testing.T) {
	type disFoodTest struct {
		arg1 int
		expected food.FoodType
	}

	var disFoodTests = []disFoodTest{
		disFoodTest{0, 0},
		disFoodTest{50, 0},
		disFoodTest{100, 0},
	}
	for _, test := range disFoodTests{
		// Check no food is returned
		if output := DisFood(test.arg1); output != test.expected {
			t.Errorf("Output %d not equal to expected %d", output, test.expected)
		}
	}
}

func TestSatisfice(t *testing.T) {
	type satisficeTest struct {
		arg1 int
		expected food.FoodType
	}

	var satisficeTests = []satisficeTest{
		satisficeTest{0, 20},
		satisficeTest{15, 20},
		satisficeTest{50, 1},
		satisficeTest{100, 1},
	}
	for _, test := range satisficeTests{
		// Check food value is 20 if hp <= 20 else 1
		if output := Satisfice(test.arg1); output != test.expected {
			t.Errorf("Output %d not equal to expected %d", output, test.expected)
		}
	}
}

func TestSatisfy(t *testing.T) {
	type satisfyTest struct {
		arg1 int
		expected food.FoodType
	}

	var satisfyTests = []satisfyTest{
		satisfyTest{0, 100},
		satisfyTest{15, 85},
		satisfyTest{50, 50},
		satisfyTest{100, 0},
	}
	for _, test := range satisfyTests{
		// Check food value returned is 100 - hp
		if output := Satisfy(test.arg1); output != test.expected {
			t.Errorf("Output %d not equal to expected %d", output, test.expected)
		}
	}
}

func TestSelectAction(t *testing.T) {
	// TODO: Implement this test
}
