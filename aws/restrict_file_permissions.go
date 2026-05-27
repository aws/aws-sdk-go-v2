package aws

// RestrictFilePermissions controls whether the SDK restricts file permissions
// on credential cache files it creates.
type RestrictFilePermissions string

const (
	// Unset.
	RestrictFilePermissionsUnset RestrictFilePermissions = ""

	// Sets file permissions to owner read/write only (0600) and directory
	// permissions to owner only (0700) when creating new cache files and
	// directories on Unix. This is the default behavior.
	RestrictFilePermissionsUserReadWrite = "user_read_write"

	// Does not set any file or directory permissions, relying on the system's
	// default umask.
	RestrictFilePermissionsUnrestricted = "unrestricted"
)
