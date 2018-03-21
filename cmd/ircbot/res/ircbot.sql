CREATE TABLE irc_channel (
	id SERIAL,
	name TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL,
	answer_commands BOOLEAN NOT NULL DEFAULT 'f',
	do_hello BOOLEAN NOT NULL DEFAULT 'f',
	PRIMARY KEY (id)
);

CREATE OR REPLACE FUNCTION register_irc_channel_config(n TEXT, p TEXT, ac BOOLEAN, hello BOOLEAN)
	RETURNS SERIAL AS $$
BEGIN
	RETURN QUERY INSERT INTO irc_channel (name, password, answer_commands, do_hello) VALUES (n, p, ac, hello) ON CONFLICT ON
		CONSTRAINT irc_channel_name_key DO UPDATE SET name = n, password = p, answer_commands = ac, do_hello = hello
	RETURNING id;
END;
$$ LANGUAGE plpgsql;