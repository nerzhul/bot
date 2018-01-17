#!/bin/sh
# PROVIDE: matterbot
# REQUIRE: DAEMON NETWORKING
# KEYWORD: shutdown

#
# Add the following line to /etc/rc.conf to enable matterbot:
#
# matterbot_enable (bool):  Set to "NO" by default.
#                               Set it to "YES" to enable matterbot
# matterbot_config (str):   Set to "" by default.
#                               Set it to matterbot configuration file
# matterbot_user (str):     Set to "matterbot" by default.
#                               Set it to user to run matterbot under
# matterbot_group (str):    Set to "matterbot" by default.
#                               Set it to group to run matterbot under

. /etc/rc.subr

name="matterbot"
rcvar="matterbot_enable"

load_rc_config $name

: ${matterbot_enable:="NO"}
: ${matterbot_config:=""}
: ${matterbot_user:="matterbot"}
: ${matterbot_group:="matterbot"}

pidfile="/var/run/${name}.pid"
procname=/usr/local/bin/matterbot
command="/usr/sbin/daemon"
command_args="-f -p ${pidfile} ${procname}"
if [ "x${matterbot_config}" != "x" ]; then
        command_args="${command_args} --config ${matterbot_config}"
fi

start_precmd="matterbot_startprecmd"

matterbot_startprecmd()
{
        if [ ! -e "${pidfile}" ]; then
                install -g ${matterbot_group} -o ${matterbot_user} -- /dev/null "${pidfile}";
        fi
}
run_rc_command $1
