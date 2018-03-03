package rabbitmq

import (
	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClient(t *testing.T) {
	assert.NotEqual(t, nil, NewClient(&logging.Logger{}, &testCfg, nil))
}

func TestClient_AddConsumerName(t *testing.T) {
	c := NewClient(&logging.Logger{}, &testCfg, nil)
	assert.NotEqual(t, nil, c)
	c.AddConsumerName("test_consumer")
	var found = false
	for _, cn := range c.consumerNames {
		if cn == "test_consumer" {
			found = true
		}
	}

	assert.Equal(t, true, found)
}

func TestClient_VerifyPublisher(t *testing.T) {
	c := NewClient(&logging.Logger{}, &testCfg, nil)
	assert.NotEqual(t, nil, c)

	assert.Equal(t, true, c.VerifyPublisher())
}

func TestClient_PublishCommand(t *testing.T) {
	c := NewClient(&logging.Logger{}, &testCfg, nil)
	assert.NotEqual(t, nil, c)

	assert.Equal(t, true, c.VerifyPublisher())

	assert.Equal(t, true, c.PublishCommand(&CommandEvent{
		Command: "help",
		User:    "unittest",
		Channel: "ci",
	}, "blackhole"))
}

func TestClient_PublishGitlabEvent(t *testing.T) {
	c := NewClient(&logging.Logger{}, &testCfg, nil)
	assert.NotEqual(t, nil, c)

	assert.Equal(t, true, c.VerifyPublisher())

	assert.Equal(t, true, c.PublishGitlabEvent(&CommandResponse{
		Message:     "event",
		MessageType: "gitlab-event",
		User:        "unittest",
		Channel:     "ci",
	}))
}
