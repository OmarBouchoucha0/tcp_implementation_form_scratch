#!/bin/sh

main=./main/*.go
output=main.o

go build -o $output $main 
