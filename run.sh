#!/bin/bash

rm -rf gotrew
go build
if [ $? -eq 0 ]; then
    ./gotrew
else
    echo "FAIL on go build"
fi
