package rabbitmq

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var cr1 = CommandResponse{
	Channel: "test",
	User:    "toto",
	Message: "testmessage",
}

const cr1TestResult = `{"channel":"test","message":"testmessage","user":"toto","message_type":""}`

func TestCommandResponse_ToJSON(t *testing.T) {
	cr1b, err := cr1.ToJSON()
	if err != nil {
		t.Errorf("Unable to convert cr1 to JSON")
	}

	if string(cr1b) != cr1TestResult {
		t.Errorf("Invalid JSON conversion for cr1. Found %s. Expected %s", cr1b, cr1TestResult)
	}
}

var ic1 = IRCCommand{
	Channel: "test",
	User:    "toto",
	Command: "join",
	Arg1:    "test",
	Arg2:    "test",
}

var ic2invalid = IRCCommand{
	Channel: "test",
	Command: "invalid",
}

const ic1TestResult = `{"command":"join","arg1":"test","arg2":"test","channel":"test","user":"toto"}`

func TestIRCCommand_ToJSON(t *testing.T) {
	ic1b, err := ic1.ToJSON()
	if err != nil {
		t.Errorf("Unable to convert cr1 to JSON")
	}

	if string(ic1b) != ic1TestResult {
		t.Errorf("Invalid JSON conversion for ic1. Found %s. Expected %s", ic1b, ic1TestResult)
	}
}

func TestIRCCommand_ToJSON_InvalidCases(t *testing.T) {
	ic2b, err := ic2invalid.ToJSON()
	assert.NotEqual(t, nil, err)
	assert.Nil(t, ic2b)
}
