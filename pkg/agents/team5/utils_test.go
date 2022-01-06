package team5

import (
	"testing"

	"github.com/SOMAS2021/SOMAS2021/pkg/messages"
)

func Test_statementsIntersect(t *testing.T) {
	type args struct {
		op1    messages.Op
		value1 int
		op2    messages.Op
		value2 int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		//op1 EQ
		{
			name: "EQ EQ true",
			args: args{
				op1:    messages.EQ,
				value1: 1,
				op2:    messages.EQ,
				value2: 1,
			},
			want: true,
		},
		{
			name: "EQ EQ false",
			args: args{
				op1:    messages.EQ,
				value1: 1,
				op2:    messages.EQ,
				value2: 2,
			},
			want: false,
		},
		{
			name: "EQ LT true",
			args: args{
				op1:    messages.EQ,
				value1: 1,
				op2:    messages.LT,
				value2: 2,
			},
			want: true,
		},
		{
			name: "EQ LT false",
			args: args{
				op1:    messages.EQ,
				value1: 1,
				op2:    messages.LT,
				value2: 1,
			},
			want: false,
		},
		{
			name: "EQ LE true",
			args: args{
				op1:    messages.EQ,
				value1: 1,
				op2:    messages.LE,
				value2: 1,
			},
			want: true,
		},
		{
			name: "EQ LE false",
			args: args{
				op1:    messages.EQ,
				value1: 1,
				op2:    messages.LE,
				value2: 0,
			},
			want: false,
		},
		{
			name: "EQ GT true",
			args: args{
				op1:    messages.EQ,
				value1: 1,
				op2:    messages.GT,
				value2: 0,
			},
			want: true,
		},
		{
			name: "EQ GT false",
			args: args{
				op1:    messages.EQ,
				value1: 1,
				op2:    messages.GT,
				value2: 1,
			},
			want: false,
		},
		{
			name: "EQ GE true",
			args: args{
				op1:    messages.EQ,
				value1: 1,
				op2:    messages.GE,
				value2: 1,
			},
			want: true,
		},
		{
			name: "EQ GE false",
			args: args{
				op1:    messages.EQ,
				value1: 1,
				op2:    messages.GE,
				value2: 2,
			},
			want: false,
		},

		// op1 LT
		{
			name: "LT EQ true",
			args: args{
				op1:    messages.LT,
				value1: 2,
				op2:    messages.EQ,
				value2: 1,
			},
			want: true,
		},
		{
			name: "LT EQ false",
			args: args{
				op1:    messages.LT,
				value1: 1,
				op2:    messages.EQ,
				value2: 1,
			},
			want: false,
		},
		{
			name: "LT LT true",
			args: args{
				op1:    messages.LT,
				value1: 1,
				op2:    messages.LT,
				value2: 1,
			},
			want: true,
		},
		{
			name: "LT LE true",
			args: args{
				op1:    messages.LT,
				value1: 1,
				op2:    messages.LE,
				value2: 1,
			},
			want: true,
		},
		{
			name: "LT GT true",
			args: args{
				op1:    messages.LT,
				value1: 3,
				op2:    messages.GT,
				value2: 1,
			},
			want: true,
		},
		{
			name: "LT GT false",
			args: args{
				op1:    messages.LT,
				value1: 2,
				op2:    messages.GT,
				value2: 1,
			},
			want: false,
		},
		{
			name: "LT GT false",
			args: args{
				op1:    messages.LT,
				value1: 1,
				op2:    messages.GT,
				value2: 1,
			},
			want: false,
		},
		{
			name: "LT GE true",
			args: args{
				op1:    messages.LT,
				value1: 2,
				op2:    messages.GE,
				value2: 1,
			},
			want: true,
		},
		{
			name: "LT GE false",
			args: args{
				op1:    messages.LT,
				value1: 1,
				op2:    messages.GE,
				value2: 1,
			},
			want: false,
		},

		// op1 LE
		{
			name: "LE EQ true",
			args: args{
				op1:    messages.LE,
				value1: 1,
				op2:    messages.EQ,
				value2: 1,
			},
			want: true,
		},
		{
			name: "LE EQ false",
			args: args{
				op1:    messages.LE,
				value1: 0,
				op2:    messages.EQ,
				value2: 1,
			},
			want: false,
		},
		{
			name: "LE LT true",
			args: args{
				op1:    messages.LE,
				value1: 1,
				op2:    messages.LT,
				value2: 1,
			},
			want: true,
		},
		{
			name: "LE LE true",
			args: args{
				op1:    messages.LE,
				value1: 1,
				op2:    messages.LE,
				value2: 1,
			},
			want: true,
		},
		{
			name: "LE GT true",
			args: args{
				op1:    messages.LE,
				value1: 2,
				op2:    messages.GT,
				value2: 1,
			},
			want: true,
		},
		{
			name: "LE GT false",
			args: args{
				op1:    messages.LE,
				value1: 1,
				op2:    messages.GT,
				value2: 1,
			},
			want: false,
		},
		{
			name: "LE GE true",
			args: args{
				op1:    messages.LE,
				value1: 1,
				op2:    messages.GE,
				value2: 1,
			},
			want: true,
		},
		{
			name: "LE GE false",
			args: args{
				op1:    messages.LE,
				value1: 0,
				op2:    messages.GE,
				value2: 1,
			},
			want: false,
		},
		// op1 GT
		{
			name: "GT EQ true",
			args: args{
				op1:    messages.GT,
				value1: 0,
				op2:    messages.EQ,
				value2: 1,
			},
			want: true,
		},
		{
			name: "GT EQ false",
			args: args{
				op1:    messages.GT,
				value1: 1,
				op2:    messages.EQ,
				value2: 1,
			},
			want: false,
		},
		{
			name: "GT LT true",
			args: args{
				op1:    messages.GT,
				value1: 0,
				op2:    messages.LT,
				value2: 2,
			},
			want: true,
		},
		{
			name: "GT LT false",
			args: args{
				op1:    messages.GT,
				value1: 1,
				op2:    messages.LT,
				value2: 1,
			},
			want: false,
		},
		{
			name: "GT LE true",
			args: args{
				op1:    messages.GT,
				value1: 0,
				op2:    messages.LE,
				value2: 1,
			},
			want: true,
		},
		{
			name: "GT LE false",
			args: args{
				op1:    messages.GT,
				value1: 1,
				op2:    messages.LE,
				value2: 1,
			},
			want: false,
		},
		{
			name: "GT GT true",
			args: args{
				op1:    messages.GT,
				value1: 1,
				op2:    messages.GT,
				value2: 1,
			},
			want: true,
		},
		{
			name: "GT GE true",
			args: args{
				op1:    messages.GT,
				value1: 1,
				op2:    messages.GE,
				value2: 1,
			},
			want: true,
		},
		// op1 GE
		{
			name: "GE EQ true",
			args: args{
				op1:    messages.GE,
				value1: 1,
				op2:    messages.EQ,
				value2: 1,
			},
			want: true,
		},
		{
			name: "GE EQ false",
			args: args{
				op1:    messages.GE,
				value1: 2,
				op2:    messages.EQ,
				value2: 1,
			},
			want: false,
		},
		{
			name: "GE LT true",
			args: args{
				op1:    messages.GE,
				value1: 1,
				op2:    messages.LT,
				value2: 2,
			},
			want: true,
		},
		{
			name: "GE LT false",
			args: args{
				op1:    messages.GE,
				value1: 1,
				op2:    messages.LT,
				value2: 1,
			},
			want: false,
		},
		{
			name: "GE LE true",
			args: args{
				op1:    messages.GE,
				value1: 1,
				op2:    messages.LE,
				value2: 1,
			},
			want: true,
		},
		{
			name: "GE LE false",
			args: args{
				op1:    messages.GE,
				value1: 2,
				op2:    messages.LE,
				value2: 1,
			},
			want: false,
		},
		{
			name: "GE GT true",
			args: args{
				op1:    messages.GE,
				value1: 1,
				op2:    messages.GT,
				value2: 1,
			},
			want: true,
		},
		{
			name: "GE GE true",
			args: args{
				op1:    messages.GE,
				value1: 1,
				op2:    messages.GE,
				value2: 1,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := statementsIntersect(tt.args.op1, tt.args.value1, tt.args.op2, tt.args.value2); got != tt.want {
				t.Errorf("statementsIntersect() = %v, want %v", got, tt.want)
			}
		})
	}
}
