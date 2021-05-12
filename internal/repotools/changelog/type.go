package changelog

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SemVerIncrement describes how a Annotation should affect a module's version.
type SemVerIncrement int

const (
	// DefaultBump indicates the the module's version should be incremented using the version selectors default behavior.
	DefaultBump SemVerIncrement = iota

	// PatchBump indicates the module's version should be incremented by a patch version bump.
	PatchBump
	// MinorBump indicates the module's version should be incremented by a minor version bump.
	MinorBump
	// ReleaseBump indicates the module version should be updated from a pre-release tag.
	ReleaseBump
)

// ChangeType describes the type of change made to a Go module.
type ChangeType int

// MarshalTOML marshals the ChangeType to a TOML string representation.
func (c ChangeType) MarshalTOML() ([]byte, error) {
	return []byte("\"" + c.String() + "\""), nil
}

// UnmarshalTOML unmarshal i, which must be a string, to it's ChangeType representaiton.
func (c *ChangeType) UnmarshalTOML(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("expect to unmarshal string, got %T", i)
	}
	*c = ParseChangeType(v)
	return nil
}

// MarshalJSON marshals the ChangeType to a JSON string representation.
func (c ChangeType) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// UnmarshalJSON unmarshal bytes, which must be a string, to it's ChangeType representaiton.
func (c *ChangeType) UnmarshalJSON(bytes []byte) error {
	var t string
	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	}
	*c = ParseChangeType(t)
	return nil
}

// ChangeType in order from least to most precedence when generating changelog summaries
const (
	UnknownChangeType ChangeType = iota
	// DependencyChangeType is a constant change type for a dependency update.
	DependencyChangeType
	// DocumentationChangeType is a constant change type for an SDK announcement.
	DocumentationChangeType
	// BugFixChangeType is a constant change type for a bug fix.
	BugFixChangeType
	// FeatureChangeType is a constant change type for a new feature.
	FeatureChangeType
	// ReleaseChangeType is a constant change type for a major version updates (from v0 => v1).
	ReleaseChangeType
	// AnnouncementChangeType is a constant change type for an SDK announcement.
	AnnouncementChangeType
)

// ParseChangeType attempts to parse the given string v into a ChangeType, returning an error if the string is invalid.
func ParseChangeType(v string) ChangeType {
	switch {
	case strings.EqualFold(FeatureChangeType.String(), v):
		return FeatureChangeType
	case strings.EqualFold(BugFixChangeType.String(), v):
		return BugFixChangeType
	case strings.EqualFold(ReleaseChangeType.String(), v):
		return ReleaseChangeType
	case strings.EqualFold(DependencyChangeType.String(), v):
		return DependencyChangeType
	case strings.EqualFold(AnnouncementChangeType.String(), v):
		return AnnouncementChangeType
	case strings.EqualFold(DocumentationChangeType.String(), v):
		return DocumentationChangeType
	default:
		return UnknownChangeType
	}
}

// ChangelogPrefix returns the CHANGELOG header the ChangeType should be grouped under.
func (c ChangeType) ChangelogPrefix() string {
	switch c {
	case FeatureChangeType:
		return "Feature"
	case BugFixChangeType:
		return "Bug Fix"
	case ReleaseChangeType:
		return "Release"
	case DependencyChangeType:
		return "Dependency Update"
	case DocumentationChangeType:
		return "Documentation"
	case AnnouncementChangeType:
		return "Announcement"
	default:
		return ""
	}
}

// VersionIncrement returns the SemVerIncrement corresponding to the given ChangeType.
func (c ChangeType) VersionIncrement() SemVerIncrement {
	switch c {
	case ReleaseChangeType:
		return ReleaseBump
	case FeatureChangeType:
		return MinorBump
	case BugFixChangeType:
		return PatchBump
	case DependencyChangeType:
		return PatchBump
	case DocumentationChangeType:
		return PatchBump
	case AnnouncementChangeType:
		fallthrough
	default:
		return DefaultBump
	}
}

// String returns a string representation of the ChangeType
func (c ChangeType) String() string {
	switch c {
	case AnnouncementChangeType:
		return "announcement"
	case ReleaseChangeType:
		return "release"
	case FeatureChangeType:
		return "feature"
	case BugFixChangeType:
		return "bugfix"
	case DocumentationChangeType:
		return "documentation"
	case DependencyChangeType:
		return "dependency"
	default:
		return ""
	}
}
