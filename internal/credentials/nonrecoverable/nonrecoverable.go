// Package nonrecoverable holds the error credential providers return when retrying a refresh won't help.
package nonrecoverable

// Failure is a refresh error that won't succeed on retry. Err is the source's error and ActionableMessage
// tells the customer how to fix it, like "run `aws sso login`".
type Failure struct {
	Err               error
	ActionableMessage string
}

func (e *Failure) Error() string {
	if e.Err == nil {
		return e.ActionableMessage
	}
	return e.ActionableMessage + ": " + e.Err.Error()
}

func (e *Failure) Unwrap() error {
	return e.Err
}
