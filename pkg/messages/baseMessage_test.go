package messages

import (
	"testing"
)

func TestMessageType_String(t *testing.T) {
	tests := []struct {
		name string
		m    MessageType
		want string
	}{
		{
			name: "AskFoodTaken",
			m:    AskFoodTaken,
			want: "AskFoodTaken",
		},
		{
			name: "AskHP",
			m:    AskHP,
			want: "AskHP",
		},
		{
			name: "AskFoodOnPlatform",
			m:    AskFoodOnPlatform,
			want: "AskFoodOnPlatform",
		},
		{
			name: "AskIntendedFoodIntake",
			m:    AskIntendedFoodIntake,
			want: "AskIntendedFoodIntake",
		},
		{
			name: "StateFoodTaken",
			m:    StateFoodTaken,
			want: "StateFoodTaken",
		},
		{
			name: "StateHP",
			m:    StateHP,
			want: "StateHP",
		},
		{
			name: "StateFoodOnPlatform",
			m:    StateFoodOnPlatform,
			want: "StateFoodOnPlatform",
		},
		{
			name: "StateIntendedFoodIntake",
			m:    StateIntendedFoodIntake,
			want: "StateIntendedFoodIntake",
		},
		{
			name: "StateIdentity",
			m:    StateIdentity,
			want: "StateIdentity",
		},
		{
			name: "StateResponse",
			m:    StateResponse,
			want: "StateResponse",
		},
		{
			name: "RequestLeaveFood",
			m:    RequestLeaveFood,
			want: "RequestLeaveFood",
		},
		{
			name: "RequestTakeFood",
			m:    RequestTakeFood,
			want: "RequestTakeFood",
		},
		{
			name: "Response",
			m:    Response,
			want: "Response",
		},
		{
			name: "UNKNOWN",
			m:    1234,
			want: "UNKNOWN",
		},
		{
			name: "zero",
			m:    0,
			want: "UNKNOWN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("MessageType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
