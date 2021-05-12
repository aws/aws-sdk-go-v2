package git

// Tag invokes git-tag is the specified path, creating the given annotated tag at the given commit with the
// provided message.
func Tag(path, tag, message string, commit string) error {
	arguments := []string{"tag", "-a", "-m", message, tag, commit}
	_, err := Git(path, arguments...)
	return err
}
