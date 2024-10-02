package aws

// RequestChecksumCalculation controls request checksum calculation workflow
type RequestChecksumCalculation int

const (
	// RequestChecksumCalculationUnset is the unset value for RequestChecksumCalculation
	RequestChecksumCalculationUnset RequestChecksumCalculation = iota

	// RequestChecksumCalculationWhenSupported indicates request checksum should be calculated
	// when the operation supports input checksums
	RequestChecksumCalculationWhenSupported

	// RequestChecksumCalculationWhenRequired indicates request checksum should be calculated
	// when user sets a checksum algorithm
	RequestChecksumCalculationWhenRequired
)

// ResponseChecksumValidation controls response checksum validation workflow
type ResponseChecksumValidation int

const (
	// ResponseChecksumValidationUnset is the unset value for ResponseChecksumValidation
	ResponseChecksumValidationUnset ResponseChecksumValidation = iota

	// ResponseChecksumValidationWhenSupported indicates response checksum should be validated
	// when the operation supports output checksums
	ResponseChecksumValidationWhenSupported

	// ResponseChecksumValidationWhenRequired indicates response checksum should be validated
	// when user enables that in validation mode cfg
	ResponseChecksumValidationWhenRequired
)
