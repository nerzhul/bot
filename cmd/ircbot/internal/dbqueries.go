package internal

const (
	// ValidationQuery query used to validate connection
	ValidationQuery = `SELECT 1`
	// LoadChannelConfigQuery select IRC configurations for channels from database
	LoadChannelConfigQuery = `SELECT name, password, answer_commands, do_hello FROM irc_channel`
)
