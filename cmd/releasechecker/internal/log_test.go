package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_InitLogger(t *testing.T) {
	assert.Equal(t, true, initLogger())
}
