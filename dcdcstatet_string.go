// Code generated by "stringer -type=DcdcStatet"; DO NOT EDIT.

package dcdcusb

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[StateOk-7]
	_ = x[StateIgnOff-8]
	_ = x[StateHardOffCountdown-16]
	_ = x[StateUnknown-255]
}

const (
	_DcdcStatet_name_0 = "StateOkStateIgnOff"
	_DcdcStatet_name_1 = "StateHardOffCountdown"
	_DcdcStatet_name_2 = "StateUnknown"
)

var (
	_DcdcStatet_index_0 = [...]uint8{0, 7, 18}
)

func (i DcdcStatet) String() string {
	switch {
	case 7 <= i && i <= 8:
		i -= 7
		return _DcdcStatet_name_0[_DcdcStatet_index_0[i]:_DcdcStatet_index_0[i+1]]
	case i == 16:
		return _DcdcStatet_name_1
	case i == 255:
		return _DcdcStatet_name_2
	default:
		return "DcdcStatet(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
