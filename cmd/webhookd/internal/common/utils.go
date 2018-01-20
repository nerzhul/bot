package common

import (
	"encoding/json"
	"io"
)

// ReadJSONRequest Helper to decode a json payload
func ReadJSONRequest(payload io.ReadCloser, decodedPayload interface{}) bool {
	decoder := json.NewDecoder(payload)
	defer payload.Close()
	if err := decoder.Decode(&decodedPayload); err != nil {
		return false
	}

	return true
}
