package slices

import "errors"

func Find[T any](ss []T, test func(T) bool) (*T, error) {
	for _, s := range ss {
		if test(s) {
			return &s, nil
		}
	}
	return nil, errors.New("element not found")
}
