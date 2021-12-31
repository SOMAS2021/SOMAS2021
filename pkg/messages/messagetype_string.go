// Code generated by "stringer -type=MessageType"; DO NOT EDIT.

package messages

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[AskFoodTaken-1]
	_ = x[AskHP-2]
	_ = x[AskFoodOnPlatform-3]
	_ = x[AskIntendedFoodIntake-4]
	_ = x[AskIdentity-5]
	_ = x[StateFoodTaken-6]
	_ = x[StateHP-7]
	_ = x[StateFoodOnPlatform-8]
	_ = x[StateIntendedFoodIntake-9]
	_ = x[StateResponse-10]
	_ = x[ProposeTreaty-11]
	_ = x[RequestLeaveFood-12]
	_ = x[RequestTakeFood-13]
	_ = x[Response-14]
	_ = x[TreatyResponse-15]
}

const _MessageType_name = "AskFoodTakenAskHPAskFoodOnPlatformAskIntendedFoodIntakeAskIdentityStateFoodTakenStateHPStateFoodOnPlatformStateIntendedFoodIntakeStateResponseProposeTreatyRequestLeaveFoodRequestTakeFoodResponseTreatyResponse"

var _MessageType_index = [...]uint8{0, 12, 17, 34, 55, 66, 80, 87, 106, 129, 142, 155, 171, 186, 194, 208}

func (i MessageType) String() string {
	i -= 1
	if i < 0 || i >= MessageType(len(_MessageType_index)-1) {
		return "MessageType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _MessageType_name[_MessageType_index[i]:_MessageType_index[i+1]]
}