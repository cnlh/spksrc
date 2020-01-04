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
service_prestart ()
{
	chown root:root "${INSTALL_DIR}/ui/index.cgi"
}
service_postinst ()
{	
	chown root:root "${INSTALL_DIR}/ui/index.cgi"	
#	${BITLBEE} install -server=1.1.1.1 -vkey=123 -type=tcp
}

service_postuninst ()
{
	${BITLBEE} uninstall
}

start_daemon ()
{
    chown root:root "${INSTALL_DIR}/ui/index.cgi"
    chmod 777 "${INSTALL_DIR}/conf/config.conf"
    var=$(cat "${INSTALL_DIR}/conf/config.conf")
    server=`echo $var | cut -d \# -f 1`
    vkey=`echo $var | cut -d \# -f 2`
    tp=`echo $var |cut -d \# -f 3`
    nohup ${BITLBEE}  -server=$server -vkey=$vkey -type=$tp &
}

stop_daemon ()
{
    killall npc
}

daemon_status ()
{
   ${BITLBEE} status
}


case $1 in
    start)
#        if daemon_status; then
#            echo ${DNAME} is already running
#            exit 0
#        else
#            echo Starting ${DNAME} ...
            start_daemon
#            exit $?
#        fi
        ;;
    stop)
#        if daemon_status; then
#            echo Stopping ${DNAME} ...
            stop_daemon
#            exit $?
#        else
#            echo ${DNAME} is not running
#            exit 0
#        fi
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
