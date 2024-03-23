#!/bin/sh
echo "building ..."
./build.sh
echo "setting permissions ..."
sudo setcap CAP_NET_ADMIN+eip main.o
./main.o &
pid=$!
echo "configuring ip adresse ..."
sudo ip addr add 192.168.0.1/24 dev tun0
sudo ip link set dev tun0 up
wait $pid
sudo ip link delete tun0
echo "deleting tun0 and exiting ..."
