#! /bin/bash
# filename: run.sh
export AK=""
export AKSECRET=""
export ENDPOINT=""
function killProsess() {
        NAME=$1

        PID=$(ps -e | grep $NAME | awk '{print $1}')

        echo "PID: $PID"
        kill -9 $PID
}

function start() {
        echo "start go-logrotate ..."
        `echo "====================================" >> ./cron.log`
        `./go-logrotate &>> ./cron.log &`
}

function stop() {
        echo "stop go-logrotate ..."
        killProsess "go-logrotate"
}

function restart() {
        echo "restart go-logrotate ..."
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
                echo "**********************"
                echo "start | restart | stop"
                echo "**********************"
                ;;
esac