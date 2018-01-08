package internal

import "encoding/base64"

func (r *commandRouter) handlerB64decode(args string, user string, channel string) *string {
	result := new(string)
	b64decoded, err := base64.StdEncoding.DecodeString(args)
	if err != nil {
		*result = "Failed to decode string: " + err.Error()
		return result
	}
	*result = "Result: " + string(b64decoded)
	return result
}

func (r *commandRouter) handlerB64encode(args string, user string, channel string) *string {
	result := new(string)
	*result = "Result: " + base64.StdEncoding.EncodeToString([]byte(args))
	return result
}
