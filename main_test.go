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

func TestSumChars(t *testing.T) {
	hash1 := "d286fa96053d0b18502c1b8ea77420c6"
	hash2 := "713de9f13e306102417c2930dc928e43"
	res1 := sumChars(hash1)
	res2 := sumChars(hash2)
	assert.NotEqual(t, res1, res2)
}
