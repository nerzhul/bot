package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIrcDB_ValidationQuery(t *testing.T) {
	assert.Equal(t, true, gIRCDB.ValidationQuery())
}

func TestIrcDB_loadIRCChannelConfigs(t *testing.T) {
	chanConfig, err := gIRCDB.loadIRCChannelConfigs()
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(chanConfig))
}
