#!/usr/bin/env bash

cd ${0%/*} #Set working dir to current dir

RUN_BINARY_NAME=""
DO_BUILD=0
DO_RUN=0


while getopts "brn:" opt; do
    case "$opt" in
    b)
        DO_BUILD=1
        ;;
    r)
        DO_RUN=1
        ;;
    n)
        RUN_BINARY_NAME="$OPTARG"
        ;;
    esac
done

if [ -z "$RUN_BINARY_NAME" ]; then   #Asume directory name as binary name if no name given
    RUN_BINARY_NAME=${PWD##*/}
fi

if [ $DO_BUILD -eq 1 ]; then
    go build -a
fi

if [ ! -f "$RUN_BINARY_NAME" ]; then
    echo "Binary $RUN_BINARY_NAME does not exist, build it with -b or use -n to specify name."
    exit 1

fi

if [ $DO_RUN -eq 1 ]; then
     pkill $RUN_BINARY_NAME
    ./$RUN_BINARY_NAME
fi




