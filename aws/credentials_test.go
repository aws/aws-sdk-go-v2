package aws

import "time"

type stubProvider struct {
	creds   Credentials
	expires time.Time
	err     error

	onInvalidate func(*stubProvider)
}

func (s *stubProvider) Retrieve() (Credentials, error) {
	creds := s.creds
	creds.Source = "stubProvider"
	creds.CanExpire = !s.expires.IsZero()
	creds.Expires = s.expires

	return creds, s.err
}

func (s *stubProvider) Invalidate() {
	s.onInvalidate(s)
}
