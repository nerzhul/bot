#!/bin/sh
# PROVIDE: slackbot
# REQUIRE: DAEMON NETWORKING
# KEYWORD: shutdown

#
# Add the following line to /etc/rc.conf to enable slackbot:
#
# slackbot_enable (bool):  Set to "NO" by default.
#                               Set it to "YES" to enable slackbot
# slackbot_config (str):   Set to "" by default.
#                               Set it to slackbot configuration file
# slackbot_user (str):     Set to "slackbot" by default.
#                               Set it to user to run slackbot under
# slackbot_group (str):    Set to "slackbot" by default.
#                               Set it to group to run slackbot under

. /etc/rc.subr

name="slackbot"
rcvar="slackbot_enable"

load_rc_config $name

: ${slackbot_enable:="NO"}
: ${slackbot_config:=""}
: ${slackbot_user:="slackbot"}
: ${slackbot_group:="slackbot"}

pidfile="/var/run/${name}.pid"
procname=/usr/local/bin/slackbot
command="/usr/sbin/daemon"
command_args="-f -p ${pidfile} ${procname}"
if [ "x${slackbot_config}" != "x" ]; then
        command_args="${command_args} --config ${slackbot_config}"
fi

start_precmd="slackbot_startprecmd"

slackbot_startprecmd()
{
        if [ ! -e "${pidfile}" ]; then
                install -g ${slackbot_group} -o ${slackbot_user} -- /dev/null "${pidfile}";
        fi
}
run_rc_command $1
