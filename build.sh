#!/bin/sh

input=./cmd/*.go
output=main.o

go build -o $output $input

