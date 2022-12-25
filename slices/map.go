package slices

func Map[T any, U any](values []T, mapper func(T) U) []U {
	var ret []U
	for _, v := range values {
		ret = append(ret, mapper(v))
	}
	return ret
}
