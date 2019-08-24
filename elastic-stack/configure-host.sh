#!/bin/sh
screen -dmS elastic_screen_session ~/Library/Containers/com.docker.docker/Data/vms/0/tty
sleep 0.5
screen -S elastic_screen_session -X "sysctl -w vm.max_map_count=262144\n"
#screen -S elastic_screen_session -X "ulimit -n 65535\n"
#screen -S elastic_screen_session -X "ulimit -u 4096\n"
sleep 0.5
screen -XS elastic_screen_session quit