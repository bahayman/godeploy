#!/bin/bash
#
# godeploy    Go Github WebHook Deployment Utility
#
# chkconfig: 345 70 30
# description: Go Github WebHook Deployment Utility
# processname: godeploy

# Source function library.
. /etc/init.d/functions

RETVAL=0

USER=nginx
PROG=godeploy
BINARY=/srv/www/godeploy/$PROG

LOCKFILE=/var/lock/subsys/$PROG

start() {
        echo -n "Starting $PROG: "
        daemon --user $USER $BINARY &
        RETVAL=$?
        [ $RETVAL -eq 0 ] && touch $LOCKFILE
        echo
        return $RETVAL
}

stop() {
        echo -n "Shutting down $PROG: "
        killproc $BINARY && success || failure
        RETVAL=$?
        [ $RETVAL -eq 0 ] && rm -f $LOCKFILE
        echo
        return $RETVAL
}

status() {
        echo -n "Checking $PROG status: "
        ss -anp | grep $PROG
        RETVAL=$?
        return $RETVAL
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    status)
        status
        ;;
    restart)
        stop
        start
        ;;
    *)
        echo "Usage: $PROG {start|stop|status|restart}"
        exit 1
        ;;
esac
exit $RETVAL