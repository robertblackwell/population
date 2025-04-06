package populations

import (
	"maps"
	"slices"
	"sort"
)

// Generic - Extracts the string keys from a value of type map[string]T
// and sorts the keys
func SortedMapKeys[T any](m map[string]T) []string {
	keys := slices.Collect(maps.Keys(m))
	sort.Strings(keys)
	return keys
}
