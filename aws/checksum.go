package aws

// RequestChecksumCalculation controls request checksum calculation workflow
type RequestChecksumCalculation string

const (
	// RequestChecksumCalculationWhenSupported indicates request checksum should be calculated if modeled
	RequestChecksumCalculationWhenSupported RequestChecksumCalculation = "whenSupported"

	// RequestChecksumCalculationWhenRequired indicates request checksum should be calculated
	// if modeled and user set an algorithm
	RequestChecksumCalculationWhenRequired = "whenRequired"
)

// ResponseChecksumValidation controls response checksum validation workflow
type ResponseChecksumValidation string

const (
	// ResponseChecksumValidationWhenSupported indicates response checksum should be validated if modeled
	ResponseChecksumValidationWhenSupported ResponseChecksumValidation = "whenSupported"

	// ResponseChecksumValidationWhenRequired indicates response checksum should be validated if modeled
	// and user enable that in vlidation mode cfg
	ResponseChecksumValidationWhenRequired = "whenRequired"
)

// RequireChecksum indicates if a checksum needs calculated/validated for a request/response
type RequireChecksum string

const (
	// RequireChecksumTrue indicates checksum should be calculated/validated
	RequireChecksumTrue RequireChecksum = "true"

	// RequireChecksumFalse indicates checksum should not be calculated/validated
	RequireChecksumFalse RequireChecksum = "false"

	// RequireChecksumPending indicates further check is needed to decide
	RequireChecksumPending RequireChecksum = "pending"
)
