#!/bin/sh
# PROVIDE: ircbot
# REQUIRE: DAEMON NETWORKING
# KEYWORD: shutdown

#
# Add the following line to /etc/rc.conf to enable ircbot:
#
# ircbot_enable (bool):  Set to "NO" by default.
#                               Set it to "YES" to enable ircbot
# ircbot_config (str):   Set to "" by default.
#                               Set it to ircbot configuration file
# ircbot_user (str):     Set to "ircbot" by default.
#                               Set it to user to run ircbot under
# ircbot_group (str):    Set to "ircbot" by default.
#                               Set it to group to run ircbot under

. /etc/rc.subr

name="ircbot"
rcvar="ircbot_enable"

load_rc_config $name

: ${ircbot_enable:="NO"}
: ${ircbot_config:=""}
: ${ircbot_user:="ircbot"}
: ${ircbot_group:="ircbot"}

pidfile="/var/run/${name}.pid"
procname=/usr/local/bin/ircbot
command="/usr/sbin/daemon"
command_args="-f -p ${pidfile} ${procname}"
if [ "x${ircbot_config}" != "x" ]; then
        command_args="${command_args} --config ${ircbot_config}"
fi

start_precmd="ircbot_startprecmd"

ircbot_startprecmd()
{
        if [ ! -e "${pidfile}" ]; then
                install -g ${ircbot_group} -o ${ircbot_user} -- /dev/null "${pidfile}";
        fi
}
run_rc_command $1
