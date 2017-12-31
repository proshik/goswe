#!/bin/bash

rm -rf gotrew
go build
if [ $? -eq 0 ]; then
    ./gotrew $1
else
    echo "FAIL on go build"
fi
