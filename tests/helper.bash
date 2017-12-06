#!/bin/bash

function start_container() {
    bin/lctn -logtostderr -root `pwd`/rootfs /hello &
    wait
}

function wait() {
    time=0
    until pid=$(get_pid)
    do
        time=$(($time + 1))
        if [ $time -gt 10 ]; then
            return 1
        fi
        sleep 1
    done
}

function cleanup() {
    kill -9 $(get_pid)
}

function get_pid() {
    echo $(pidof hello)
}
