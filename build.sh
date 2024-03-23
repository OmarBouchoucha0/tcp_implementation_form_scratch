#!/bin/sh

go build -o tuntap ./cmd/tuntap.go
sudo setcap CAP_NET_ADMIN+eip tuntap
./tuntap
sudo ip addr add 192.168.0.1/24 dev tun0
sudo ip link set dev tun0 up

