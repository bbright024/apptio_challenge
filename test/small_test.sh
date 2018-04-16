#!/bin/bash



HOME_DIR="/home/bbright/go/src/apptio"
LOGSERVER_DIR="$HOME_DIR/logserver"
CONFIGS_DIR="$HOME_DIR/configs"
TEST_DIR="$HOME_DIR/test/"

cd $LOGSERVER_DIR
go test

if [ ! $? -eq 0 ]
then
    exit 1
fi

cd $CONFIGS_DIR
go test

if [ ! $? -eq 0 ]
then
    exit 1
fi

echo "All unit tests passed"
