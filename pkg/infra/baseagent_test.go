package infra

import (
	"reflect"
	"testing"

	"github.com/SOMAS2021/SOMAS2021/pkg/utils/globalTypes/food"
)

func TestBase_CurrPlatFood(t *testing.T) {
	type fields struct {
		floor     int
		platFloor int
	}
	tests := []struct {
		name   string
		fields fields
		want   food.FoodType
	}{
		{
			name: "same floor",
			fields: fields{
				floor:     1,
				platFloor: 1,
			},
			want: 100,
		},
		{
			name: "next floor",
			fields: fields{
				floor:     1,
				platFloor: 2,
			},
			want: 100,
		},
		{
			name: "wrong floor",
			fields: fields{
				floor:     2,
				platFloor: 1,
			},
			want: -1,
		},
		{
			name: "very wrong floor",
			fields: fields{
				floor:     100,
				platFloor: 1,
			},
			want: -1,
		},
		{
			name: "just missed floor",
			fields: fields{
				floor:     1,
				platFloor: 3,
			},
			want: -1,
		},
		{
			name: "very wrong floor",
			fields: fields{
				floor:     1,
				platFloor: 100,
			},
			want: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tower := &Tower{
				currPlatFood:  100,
				currPlatFloor: tt.fields.platFloor,
			}
			a := &Base{
				floor: tt.fields.floor,
				tower: tower,
			}
			if got := a.CurrPlatFood(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Base.CurrPlatFood() = %v, want %v", got, tt.want)
			}
		})
	}
}
