package git

// Commit commits the staged contents in the specified repository path using the provided message.
func Commit(path, message string) error {
	arguments := []string{"commit", "-m", message}
	_, err := Git(path, arguments...)
	return err
}
