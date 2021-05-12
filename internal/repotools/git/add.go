package git

// Add invokes git-add in the specified path passing the provided arguments.
func Add(path string, args ...string) error {
	arguments := []string{"add"}

	if len(args) > 0 {
		arguments = append(arguments, args...)
	}

	_, err := Git(path, arguments...)
	return err
}
