#! /bin/bash
# filename: run.sh

function killProsess() {
        NAME=$1

        PID=$(ps -e | grep $NAME | awk '{print $1}')

        echo "PID: $PID"
        kill -9 $PID
}

function start() {
        echo "start go-cron ..."
        `./logrotate &>> ./cron.log &`
}

function stop() {
        echo "stop go-cron ..."
        killProsess "logrotate"
}

function restart() {
        echo "restart go-cron ..."
        stop
        start
}

case "$1" in
        start )
                start
                ;;
        stop )
                stop
                ;;
        restart )
                restart
                ;;
        * )
                echo "****************"
                echo "start | restart | stop"
                echo "****************"
                ;;
esac
