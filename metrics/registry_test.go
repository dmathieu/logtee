package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistryProvider(t *testing.T) {
	p := RegistryProvider()
	p2 := RegistryProvider()

	assert.Equal(t, p, p2)
}
