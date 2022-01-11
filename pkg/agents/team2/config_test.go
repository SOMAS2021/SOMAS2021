package team2

import (
	"testing"
)

func TestInitTable(t *testing.T) {
	type initTableTest struct {
		arg1, arg2, expected1, expected2 int
	}

	var initTableTests = []initTableTest{
		initTableTest{3,3,3,3},
		initTableTest{4,50,4,50},
		initTableTest{999,10,999,10},
	}
	// Check dimensions of table are correct
	for _, test := range initTableTests{
		output := InitTable(test.arg1, test.arg2)
		if len(output) != test.expected1 {
			t.Errorf("Length %d not equal to expected %d", len(output), test.expected1)
		} else if len(output[0]) != test.expected2 {
			t.Errorf("Length %d not equal to expected %d", len(output[0]), test.expected2)
		}
	}
}

func TestNew(t *testing.T) {
	// TODO: Consider if testing is appropriate for this and implement
}

func TestRun(t *testing.T) {
	// TODO: Consider if testing is appropriate for this and implement
}
