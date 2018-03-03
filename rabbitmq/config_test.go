package rabbitmq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testCfg = Config{
	URL:           "amqp://guest:guest@rabbitmq/",
	EventExchange: "events/unittests",
	Consumers: map[string]Consumer{
		"test": {
			ConsumerID: "cid",
		},
	},
}

func TestConfig_GetConsumer(t *testing.T) {
	c := testCfg.GetConsumer("test")
	assert.NotEqual(t, nil, c)
	assert.Equal(t, "cid", c.ConsumerID)
}

func TestConfig_GetConsumerNil(t *testing.T) {
	assert.Nil(t, testCfg.GetConsumer("unknown"))
}
