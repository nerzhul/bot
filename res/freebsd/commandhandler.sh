#!/bin/sh
# PROVIDE: commandhandler
# REQUIRE: DAEMON NETWORKING
# KEYWORD: shutdown

#
# Add the following line to /etc/rc.conf to enable commandhandler:
#
# commandhandler_enable (bool):  Set to "NO" by default.
#                               Set it to "YES" to enable commandhandler
# commandhandler_config (str):   Set to "" by default.
#                               Set it to commandhandler configuration file
# commandhandler_user (str):     Set to "commandhandler" by default.
#                               Set it to user to run commandhandler under
# commandhandler_group (str):    Set to "commandhandler" by default.
#                               Set it to group to run commandhandler under

. /etc/rc.subr

name="commandhandler"
rcvar="commandhandler_enable"

load_rc_config $name

: ${commandhandler_enable:="NO"}
: ${commandhandler_config:=""}
: ${commandhandler_user:="commandhandler"}
: ${commandhandler_group:="commandhandler"}

pidfile="/var/run/${name}.pid"
procname=/usr/local/bin/commandhandler
command="/usr/sbin/daemon"
command_args="-f -p ${pidfile} ${procname}"
if [ "x${commandhandler_config}" != "x" ]; then
        command_args="${command_args} --config ${commandhandler_config}"
fi

start_precmd="commandhandler_startprecmd"

commandhandler_startprecmd()
{
        if [ ! -e "${pidfile}" ]; then
                install -g ${commandhandler_group} -o ${commandhandler_user} -- /dev/null "${pidfile}";
        fi
}
run_rc_command $1
