package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilesEqual(t *testing.T) {
	b, _ := ioutil.ReadFile("docker-compose.yaml")
	b2, _ := ioutil.ReadFile("assets/.gitkeep")
	res := filesEqual(&b, &b2)
	assert.Equal(t, false, res)
	res = filesEqual(&b, &b)
	assert.Equal(t, true, res)
}
