package enhancedclient

type SchemaOptions struct {
	ErrorOnMissingField *bool

	// FallbackErrors will be ignored by default or if set to *true
	IgnoreFallbackErrors *bool
}
