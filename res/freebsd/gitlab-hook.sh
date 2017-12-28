#!/bin/sh
# PROVIDE: gitlab-hookd
# REQUIRE: DAEMON NETWORKING
# KEYWORD: shutdown

#
# Add the following line to /etc/rc.conf to enable gitlab_hookd:
#
# gitlab_hookd_enable (bool):  Set to "NO" by default.
#                               Set it to "YES" to enable gitlab-hookd
# gitlab_hookd_config (str):   Set to "" by default.
#                               Set it to gitlab-hookd configuration file
# gitlab_hookd_user (str):     Set to "gitlab-hook" by default.
#                               Set it to user to run gitlab-hookd under
# gitlab_hookd_group (str):    Set to "gitlab-hook" by default.
#                               Set it to group to run gitlab-hookd under

. /etc/rc.subr

name="gitlab_hookd"
rcvar="gitlab_hookd_enable"

load_rc_config $name

: ${gitlab_hookd_enable:="NO"}
: ${gitlab_hookd_config:=""}
: ${gitlab_hookd_user:="gitlab-hook"}
: ${gitlab_hookd_group:="gitlab-hook"}

pidfile="/var/run/${name}.pid"
procname=/usr/local/bin/gitlab-hookd
command="/usr/sbin/daemon"
command_args="-f -p ${pidfile} ${procname}"
if [ "x${gitlab_hookd_config}" != "x" ]; then
        command_args="${command_args} --config ${gitlab_hookd_config}"
fi

start_precmd="gitlab_hookd_startprecmd"

gitlab_hookd_startprecmd()
{
        if [ ! -e "${pidfile}" ]; then
                install -g ${gitlab_hookd_group} -o ${gitlab_hookd_user} -- /dev/null "${pidfile}";
        fi
}
run_rc_command $1
