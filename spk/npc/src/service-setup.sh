#!/bin/sh

# Package


PACKAGE="npc"
DNAME="npc"

# Others
INSTALL_DIR="/usr/local/${PACKAGE}"
BITLBEE="${INSTALL_DIR}/bin/npc"
CFG_FILE="${INSTALL_DIR}/var/npc.conf"
PID_FILE="/tmp/npc.pid"

SC_USER="root"
LEGACY_USER="root"
#USER="$([ "${BUILDNUMBER}" -ge "7321" ] && echo -n ${SC_USER} || echo -n ${LEGACY_USER})"

service_postinst ()
{
	echo 
 #   ln -s ${SYNOPKG_PKGDEST}/bin/tmux-utf8 /usr/local/bin/tmux
}

service_postuninst ()
{
	echo 
 #   rm /usr/local/bin/tmux
}

start_daemon ()
{
    chown root:root "${INSTALL_DIR}/ui/index.cgi"
    chmod 777 "${INSTALL_DIR}/conf/config.conf"
    var=$(cat "${INSTALL_DIR}/conf/config.conf")
    server=`echo $var | cut -d \# -f 1`
    vkey=`echo $var | cut -d \# -f 2`
    tp=`echo $var |cut -d \# -f 3`
    ${BITLBEE} start -server=$server -vkey=$vkey -type=$tp
}

stop_daemon ()
{
    kill `cat ${PID_FILE}`
    wait_for_status 1 20 || kill -9 `cat ${PID_FILE}`
    rm -f ${PID_FILE}
}

daemon_status ()
{
    if [ -f ${PID_FILE} ] && kill -0 `cat ${PID_FILE}` > /dev/null 2>&1; then
        return
    fi
    rm -f ${PID_FILE}
    return 1
}

wait_for_status ()
{
    counter=$2
    while [ ${counter} -gt 0 ]; do
        daemon_status
        [ $? -eq $1 ] && return
        let counter=counter-1
        sleep 1
    done
    return 1
}


case $1 in
    start)
        if daemon_status; then
            echo ${DNAME} is already running
            exit 0
        else
            echo Starting ${DNAME} ...
            start_daemon
            exit $?
        fi
        ;;
    stop)
        if daemon_status; then
            echo Stopping ${DNAME} ...
            stop_daemon
            exit $?
        else
            echo ${DNAME} is not running
            exit 0
        fi
        ;;
    restart)
        stop_daemon
        start_daemon
        exit $?
        ;;
    status)
        if daemon_status; then
            echo ${DNAME} is running
            exit 0
        else
            echo ${DNAME} is not running
            exit 1
        fi
        ;;
    *)
        exit 0
        ;;
esac
