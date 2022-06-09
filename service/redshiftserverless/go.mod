module github.com/aws/aws-sdk-go-v2/service/redshiftserverless

go 1.15

retract (
	// Retract latest version of the client since module is not functional.
	v1.0.1
	// API client was incorrectly released, and is not functional.
	v1.0.0
)
