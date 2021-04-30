package repotools

import "sort"

// AppendIfNotPresent appends value to the slice x if not present. The slice must be sorted in ascending order.
func AppendIfNotPresent(x []string, value string) []string {
	i := sort.SearchStrings(x, value)
	if i < len(x) && x[i] == value {
		return x
	}
	x = append(x, "")
	copy(x[i+1:], x[i:])
	x[i] = value
	return x
}
