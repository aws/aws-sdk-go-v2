package mode

import "fmt"

// AIDMode switches on/off the account ID based endpoint routing
type AIDMode string

// enums of valid AIDMode
const (
	Preferred AIDMode = "preferred"
	Required  AIDMode = "required"
	Disabled  AIDMode = "disabled"
)

// SetFromString converts config string to corresponding AIDMode or reports error if the value is not enumerated
func (mode *AIDMode) SetFromString(s string) error {
	switch {
	case s == "preferred":
		*mode = Preferred
	case s == "required":
		*mode = Required
	case s == "disabled":
		*mode = Disabled
	default:
		return fmt.Errorf("unknown account id mode, must be preferred/required/disabled")
	}

	return nil
}
