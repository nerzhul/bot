package internal

const (
	// LoadChannelConfigQuery select IRC configurations for channels from database
	LoadChannelConfigQuery = `SELECT name, password, answer_commands, do_hello FROM irc_channel`
	// RegisterChannelConfigQuery register configuration for the specified channel
	RegisterChannelConfigQuery = `SELECT register_irc_channel_config($1, $2, 'f', 'f')`
)
