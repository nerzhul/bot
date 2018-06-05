package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_loadDefaultConfiguration(t *testing.T) {
	var c config
	assert.Equal(t, true, c.loadDefaultConfiguration())
}

// @TODO: load real configuration
