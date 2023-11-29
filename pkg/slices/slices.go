package slices

// Interfaces converts a slice of type T to a slice of interfaces.
func Interfaces[T any](vs []T) []interface{} {
	is := make([]interface{}, len(vs))
	for i, v := range vs {
		is[i] = v
	}

	return is
}
