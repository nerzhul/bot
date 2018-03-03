package rabbitmq

import (
	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClient(t *testing.T) {
	assert.NotEqual(t, nil, NewClient(&logging.Logger{}, &testCfg, nil))
}
