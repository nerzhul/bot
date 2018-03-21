package internal

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIrcDB_ValidationQuery(t *testing.T) {
	assert.Equal(t, true, gIRCDB.ValidationQuery())
}

func TestIrcDB_loadIRCChannelConfigs(t *testing.T) {
	chanConfig, err := gIRCDB.loadIRCChannelConfigs()
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(chanConfig))
}

func TestIrcDB_SaveIRCChannelConfig(t *testing.T) {
	err := gIRCDB.SaveIRCChannelConfig(fmt.Sprintf("#unittests-%d", time.Now().Unix()), "")
	assert.Nil(t, err)

	upsertTime := time.Now().Unix() + 1
	err = gIRCDB.SaveIRCChannelConfig(fmt.Sprintf("#unittests-%d", upsertTime), "pwd")
	assert.Nil(t, err)

	// Do it twice
	err = gIRCDB.SaveIRCChannelConfig(fmt.Sprintf("#unittests-%d", upsertTime), "")
	assert.Nil(t, err)
}
