#!/usr/bin/env bats

load helper

@test "test launching hello container" {
    pid=$(get_pid)
    sudo nsenter -t $pid -n ip link set lo up
    message="$(sudo nsenter -t $pid -n curl http://127.0.0.1:80)"
    [ "$message" = "Hello World from Go in minimal Docker container" ]
}

@test "test uid gid mapping" {
    pid=$(get_pid)
    map_res="$(cat /proc/$pid/uid_map)"
    [ "$map_res" = "         0       1000        100" ]
}

@test "test creating namespaces" {
    pid=$(get_pid)
    for ns in $(sudo ls /proc/1/ns/); do
        if [[ $ns != "cgroup" ]]; then
            echo $(readlink /proc/1/ns/$ns) $(sudo readlink /proc/$pid/ns/$ns)
            [ "$(sudo readlink /proc/1/ns/$ns)" != "$(sudo readlink /proc/$pid/ns/$ns)" ]
        fi
    done
}

@test "test creating devices" {
    pid=$(get_pid)
    sudo nsenter -t $pid -m echo 1 > /dev/null
    message="$(sudo nsenter -t $pid -m cat /dev/null)"
    [ "$message" = "" ]
    message="$(sudo nsenter -t $pid -m cat /proc/mounts | grep null)"
    [ "$message" = "" ]
}
