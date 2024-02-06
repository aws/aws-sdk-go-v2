package aws

// AccountIDEndpointMode switches on/off the account ID based endpoint routing
type AccountIDEndpointMode string

// enums of valid AccountIDEndpointMode
const (
	AccountIDEndpointModeUnset     AccountIDEndpointMode = ""
	AccountIDEndpointModePreferred AccountIDEndpointMode = "preferred"
	AccountIDEndpointModeRequired  AccountIDEndpointMode = "required"
	AccountIDEndpointModeDisabled  AccountIDEndpointMode = "disabled"
)
