package rabbitmq

import "testing"

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
