package shareddefaults

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// SharedCredentialsFilename returns the SDK's default file path
// for the shared credentials file.
//
// Builds the shared config file path based on the OS's platform.
//
//   - Linux/Unix: $HOME/.aws/credentials
//   - Windows: %USERPROFILE%\.aws\credentials
func SharedCredentialsFilename() string {
	return filepath.Join(UserHomeDir(), ".aws", "credentials")
}

// SharedConfigFilename returns the SDK's default file path for
// the shared config file.
//
// Builds the shared config file path based on the OS's platform.
//
//   - Linux/Unix: $HOME/.aws/config
//   - Windows: %USERPROFILE%\.aws\config
func SharedConfigFilename() string {
	return filepath.Join(UserHomeDir(), ".aws", "config")
}

// ExpandHomePath expands a leading ~ in the path to the user's home
// directory. This is necessary because the Go os package does not
// perform shell-style tilde expansion, unlike Python's
// os.path.expanduser which is used by botocore.
//
// Per the AWS SDKs and Tools Reference Guide, ~ followed by / (or the
// platform-specific path separator) at the start of a file path should
// resolve to the home directory:
// https://docs.aws.amazon.com/sdkref/latest/guide/file-location.html
func ExpandHomePath(path string) string {
	if path == "~" {
		return UserHomeDir()
	}
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~"+string(filepath.Separator)) {
		return filepath.Join(UserHomeDir(), path[2:])
	}
	return path
}

// UserHomeDir returns the home directory for the user the process is
// running under.
func UserHomeDir() string {
	// Ignore errors since we only care about Windows and *nix.
	home, _ := os.UserHomeDir()

	if len(home) > 0 {
		return home
	}

	currUser, _ := user.Current()
	if currUser != nil {
		home = currUser.HomeDir
	}

	return home
}
