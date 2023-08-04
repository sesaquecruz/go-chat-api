package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldSayHelloWorld(t *testing.T) {
	assert.Equal(t, "Hello, World!", greeting())
}
