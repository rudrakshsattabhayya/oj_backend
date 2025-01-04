package helpers

func Difference[T comparable](slice1, slice2 []T) []T {
	// Create a map to store elements of slice2
	elements := make(map[T]bool)
	for _, item := range slice2 {
		elements[item] = true
	}

	// Find elements in slice1 that are not in slice2
	var diff []T
	for _, item := range slice1 {
		if !elements[item] {
			diff = append(diff, item)
		}
	}
	return diff
}
