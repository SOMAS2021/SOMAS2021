package utilFunctions

import (
	"testing"

	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/agent"
)

func TestMinInt(t *testing.T) {
	type args struct {
		vars []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "zero",
			args: args{vars: []int{0, 5, 2, 23, 46, 5}},
			want: 0,
		},
		{
			name: "negative",
			args: args{vars: []int{0, 5, 0, -23, -46, 5}},
			want: -46,
		},
		{
			name: "single",
			args: args{vars: []int{2}},
			want: 2,
		},
		{
			name: "most negative",
			args: args{vars: []int{-9223372036854775808, -9223372036854775}},
			want: -9223372036854775808,
		},
		{
			name: "most positive",
			args: args{vars: []int{9223372036854775807, 0}},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MinInt(tt.args.vars...); got != tt.want {
				t.Errorf("MinInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSum(t *testing.T) {
	type args struct {
		input map[agent.AgentType]int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "empty",
			args: args{input: make(map[agent.AgentType]int)},
			want: 0,
		},
		{
			name: "zero",
			args: args{
				input: map[agent.AgentType]int{
					agent.Team1Agent1: 0,
				},
			},
			want: 0,
		},
		{
			name: "single",
			args: args{
				input: map[agent.AgentType]int{
					agent.Team1Agent1: 53,
				},
			},
			want: 53,
		},
		{
			name: "triangular",
			args: args{
				input: map[agent.AgentType]int{
					agent.Team1Agent1: 1,
					agent.Team1Agent2: 2,
					agent.Team2:       3,
					agent.Team3:       4,
					agent.Team4:       5,
					agent.Team5:       6,
					agent.Team6:       7,
					agent.Team7:       8,
					agent.RandomAgent: 9,
				},
			},
			want: 45,
		},
		{
			name: "negatives",
			args: args{
				input: map[agent.AgentType]int{
					agent.Team1Agent1: -1,
					agent.Team1Agent2: -2,
					agent.Team2:       -3,
					agent.Team3:       -4,
					agent.Team4:       -5,
					agent.Team5:       -6,
					agent.Team6:       -7,
					agent.Team7:       -8,
					agent.RandomAgent: -9,
				},
			},
			want: -45,
		},
		{
			name: "pos and neg",
			args: args{
				input: map[agent.AgentType]int{
					agent.Team1Agent1: -1,
					agent.Team1Agent2: 2,
					agent.Team2:       -3,
					agent.Team3:       4,
					agent.Team4:       -5,
					agent.Team5:       6,
					agent.Team6:       -7,
					agent.Team7:       8,
					agent.RandomAgent: -9,
				},
			},
			want: -5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sum(tt.args.input); got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}
