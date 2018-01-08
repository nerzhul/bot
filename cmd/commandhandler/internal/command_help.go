package internal

import "strings"

func (r *commandRouter) handlerHelp(args string, user string, channel string) *string {
	result := new(string)
	*result += "Available commands: " + strings.Join(r.commandList, ", ")
	return result
}
