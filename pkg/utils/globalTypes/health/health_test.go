package health

import (
	"reflect"
	"testing"

	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

func TestFoodRequired(t *testing.T) {
	healthInfo := &HealthInfo{
		MaxHP:          100,
		WeakLevel:      10,
		Width:          48,
		Tau:            15,
		HPReqCToW:      5,
		HPCritical:     3,
		MaxDayCritical: 3,
		HPLossBase:     5,
		HPLossSlope:    0.2,
		maxPlatFood:    100,
	}
	type args struct {
		currentHP int
		goalHP    int
	}
	tests := []struct {
		name string
		args args
		want food.FoodType
	}{
		{
			name: "maintain 100",
			args: args{
				currentHP: 100,
				goalHP:    100,
			},
			want: 13,
		},
		{
			name: "90 to 100",
			args: args{
				currentHP: 90,
				goalHP:    100,
			},
			want: 24,
		},
		{
			name: "100 to 90",
			args: args{
				currentHP: 100,
				goalHP:    90,
			},
			want: 6,
		},
		{
			name: "big jump",
			args: args{
				currentHP: 15,
				goalHP:    100,
			},
			want: healthInfo.maxPlatFood,
		},
		{
			name: "big loss",
			args: args{
				currentHP: 100,
				goalHP:    10,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FoodRequired(tt.args.currentHP, tt.args.goalHP, healthInfo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FoodRequired() = %v, want %v", got, tt.want)
			}
		})
	}
}
