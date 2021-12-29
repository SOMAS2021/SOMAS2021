package utilFunctions

import (
	"testing"
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
		inputList []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "empty",
			args: args{inputList: []int{}},
			want: 0,
		},
		{
			name: "single",
			args: args{inputList: []int{53}},
			want: 53,
		},
		{
			name: "empty",
			args: args{inputList: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
			want: 55,
		},
		{
			name: "negatives",
			args: args{inputList: []int{-1, -2, -3, -4, -5, -6, -7, -8, -9, -10}},
			want: -55,
		},
		{
			name: "pos and neg",
			args: args{inputList: []int{-1, 2, -3, 4, -5, 6, -7, 8, -9, 10}},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sum(tt.args.inputList); got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}
