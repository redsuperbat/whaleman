package slices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	slice := []string{"Hello", "World", "way!", "rofl"}
	expected := []int{72, 87, 119, 114}
	result := Map(slice, func(s string) int {
		return int(s[0])
	})
	assert.Equal(t, result, expected)
	t.Log(result)
}
