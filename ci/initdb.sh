#! /bin/bash

SQL_PATH_IRCBOT_DB="cmd/ircbot/res/"
PGPASSWORD=${POSTGRES_PASSWORD} psql -h postgres -U ${POSTGRES_USER} ${POSTGRES_DB} < ${SQL_PATH_IRCBOT_DB}/ircbot.sql