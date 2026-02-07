package v4

import "time"

// SigningTime provides a wrapper around a time.Time which provides cached values for SigV4 signing.
type SigningTime struct {
	time.Time
	timeFormat       string
	shortTimeFormat  string
	noZoneTimeFormat string
}

// NewSigningTime creates a new SigningTime given a time.Time
func NewSigningTime(t time.Time) SigningTime {
	return SigningTime{
		Time: t,
	}
}

// TimeFormat provides a time formatted in the X-Amz-Date format.
func (m *SigningTime) TimeFormat() string {
	return m.format(&m.timeFormat, TimeFormat)
}

// ShortTimeFormat provides a time formatted of 20060102.
func (m *SigningTime) ShortTimeFormat() string {
	return m.format(&m.shortTimeFormat, ShortTimeFormat)
}

// NoZoneTimeFormat provides a time formatted X-Amz-Date format (without Z suffix).
func (m *SigningTime) NoZoneTimeFormat() string {
	return m.format(&m.noZoneTimeFormat, NoZoneTimeFormat)
}

func (m *SigningTime) format(target *string, format string) string {
	if len(*target) > 0 {
		return *target
	}
	v := m.Time.Format(format)
	*target = v
	return v
}
