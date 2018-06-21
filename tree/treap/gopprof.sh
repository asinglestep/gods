#!/bin/bash

if [[ $1 == "" ]];then
    echo "输入测试函数名"
    exit
fi

go test -gcflags="-N -l" -bench=$1 -run=^$ -cpuprofile=cpu.out
go tool pprof -png cpu.out  > $1_cpupprof.png
rm *.out