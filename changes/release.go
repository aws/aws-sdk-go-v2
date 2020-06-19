package changes

type VersionBump struct {
	Module string
	From   string
	To     string
}

type Release struct {
	Id            string
	SchemaVersion string
	VersionBumps  []VersionBump
	Changes       []*Change
}

func (r *Release) SetSchemaVersion(version string) {
	r.SchemaVersion = version
}
