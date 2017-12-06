#!/usr/bin/env bats

load helper

@test "test launching hello container" {
    start_container
    pid=$(get_pid)
    nsenter -t $pid -n ip link set lo up
    message="$(nsenter -t $pid -n curl http://127.0.0.1:80)"
    [ "$message" = "Hello World from Go in minimal Docker container" ]
    cleanup
}
