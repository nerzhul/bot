CREATE TABLE irc_channel (
	id SERIAL,
	name TEXT NOT NULL,
	password TEXT NOT NULL,
	answer_commands BOOLEAN NOT NULL DEFAULT 'f',
	do_hello BOOLEAN NOT NULL DEFAULT 'f',
	PRIMARY KEY (id)
);

CREATE UNIQUE INDEX irc_channel_uidx_name ON irc_channel(name);