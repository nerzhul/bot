#!/bin/sh
# PROVIDE: twitterbot
# REQUIRE: DAEMON NETWORKING
# KEYWORD: shutdown

#
# Add the following line to /etc/rc.conf to enable twitterbot:
#
# twitterbot_enable (bool):  Set to "NO" by default.
#                               Set it to "YES" to enable twitterbot
# twitterbot_config (str):   Set to "" by default.
#                               Set it to twitterbot configuration file
# twitterbot_user (str):     Set to "twitterbot" by default.
#                               Set it to user to run twitterbot under
# twitterbot_group (str):    Set to "twitterbot" by default.
#                               Set it to group to run twitterbot under

. /etc/rc.subr

name="twitterbot"
rcvar="twitterbot_enable"

load_rc_config $name

: ${twitterbot_enable:="NO"}
: ${twitterbot_config:=""}
: ${twitterbot_user:="twitterbot"}
: ${twitterbot_group:="twitterbot"}

pidfile="/var/run/${name}.pid"
procname=/usr/local/bin/twitterbot
command="/usr/sbin/daemon"
command_args="-f -p ${pidfile} ${procname}"
if [ "x${twitterbot_config}" != "x" ]; then
        command_args="${command_args} --config ${twitterbot_config}"
fi

start_precmd="twitterbot_startprecmd"

twitterbot_startprecmd()
{
        if [ ! -e "${pidfile}" ]; then
                install -g ${twitterbot_group} -o ${twitterbot_user} -- /dev/null "${pidfile}";
        fi
}
run_rc_command $1
