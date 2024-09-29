package aws

// RequestChecksumCalculation controls request checksum calculation workflow
type RequestChecksumCalculation int

const (
	// RequestChecksumCalculationWhenSupported indicates request checksum should be calculated if
	// client operation model has request checksum trait
	RequestChecksumCalculationWhenSupported RequestChecksumCalculation = 1

	// RequestChecksumCalculationWhenRequired indicates request checksum should be calculated
	// if modeled and user set an algorithm
	RequestChecksumCalculationWhenRequired = 2
)

// ResponseChecksumValidation controls response checksum validation workflow
type ResponseChecksumValidation int

const (
	// ResponseChecksumValidationWhenSupported indicates response checksum should be validated if modeled
	ResponseChecksumValidationWhenSupported ResponseChecksumValidation = 1

	// ResponseChecksumValidationWhenRequired indicates response checksum should be validated if modeled
	// and user enable that in vlidation mode cfg
	ResponseChecksumValidationWhenRequired = 2
)
