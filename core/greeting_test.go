package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetGreeting(t *testing.T) {
	assert.Equal(t, "Hello world!", GetGreeting())
}
