#!/bin/bash

rm -rf goswe
go build
if [ $? -eq 0 ]; then
    ./goswe
else
    echo "FAIL on go build"
fi
