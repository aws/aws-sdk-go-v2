package aws

// RequestChecksumCalculation controls request checksum calculation workflow
type RequestChecksumCalculation int

const (
	// RequestChecksumCalculationWhenSupported indicates request checksum should be calculated
	// when the operation supports input checksums
	RequestChecksumCalculationWhenSupported RequestChecksumCalculation = iota

	// RequestChecksumCalculationWhenRequired indicates request checksum should be calculated
	// when user sets a checksum algorithm
	RequestChecksumCalculationWhenRequired
)

// ResponseChecksumValidation controls response checksum validation workflow
type ResponseChecksumValidation int

const (
	// ResponseChecksumValidationWhenSupported indicates response checksum should be validated
	// when the operation supports output checksums
	ResponseChecksumValidationWhenSupported ResponseChecksumValidation = iota

	// ResponseChecksumValidationWhenRequired indicates response checksum should be validated
	// when user enables that in validation mode cfg
	ResponseChecksumValidationWhenRequired
)
