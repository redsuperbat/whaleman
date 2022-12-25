package slices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	slice := []string{"Hello", "World", "", "rofl"}
	filteredSlice := []string{"Hello", "World", "rofl"}

	assert.Equal(t, Filter(slice, func(s string) bool {
		return s != ""
	}), filteredSlice)
}
