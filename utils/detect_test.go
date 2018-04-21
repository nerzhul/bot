package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestIsInDocker(t *testing.T) {
	inDocker := IsInDocker()
	_, err := os.Stat("/.dockerenv")
	assert.Equal(t, err == nil, inDocker)
	assert.Equal(t, err != nil && os.IsNotExist(err), !inDocker)
	assert.Equal(t, err != nil && !os.IsNotExist(err), inDocker)
}
