#!/bin/bash

# This script tests these steps:
#   a) compile and run the logserver
#   b) use curl to read the test mainapp.log file
#   c) compare the results file to the actual file
#        - because the log file is parsed by the server,
#          diff won't be completely accurate.  Therefore
#          this comparison can't guarantee correct log results,
#          but it can guarantee that a file was accessed.

LOGSERVER_DIR="/home/bbright/go/src/apptio/logserver"
HOME_DIR=`pwd`
LOGSERVER_CONF="conf.json"
TEST_DIR="$HOME_DIR/test"
TEST_FILE="mainapp.log"
TEST_STR="4/15/2018, This is a test string that should be downloaded "
TEST_OUTPUT="results.txt"
TIMEOUT=5
#
# download_proxy - download a file from the origin server via the proxy
# usage: download_proxy <testdir> <filename> <origin_url> <proxy_url>
#
function download_log {
    cd $1
    curl --max-time ${TIMEOUT} --silent --proxy $4 --output $2 $3
    (( $? == 28 )) && echo "Error: Fetch timed out after ${TIMEOUT} seconds"
    cd $HOME_DIR
}



echo $TEST_STR > ./$TEST_FILE


cd $LOGSERVER_DIR
go build logserver.go

if [ ! -x ./logserver ]
then 
    echo "Error: logserver not found or not an executable file."
    exit
fi

cd $HOME_DIR

if [ ! -d ${TEST_DIR} ]
then
    mkdir ${TEST_DIR}
fi

mv ./$TEST_FILE $TEST_DIR/
cp $LOGSERVER_DIR/$LOGSERVER_CONF $TEST_DIR/
mv $LOGSERVER_DIR/logserver $TEST_DIR/
cp $LOGSERVER_DIR/logserver.log $TEST_DIR/
touch $TEST_DIR/logserver.log
# Run the logserver
cd $TEST_DIR
echo "Starting logserver with conf file:"
jq . ./$LOGSERVER_CONF
./logserver &

echo "Attempting to read mainapp.log from logserver:"

curl --max-time ${TIMEOUT} --silent --output ${TEST_OUTPUT} localhost:8888/read
(( $? == 28 )) && echo "Error: Fetch timed out after ${TIMEOUT} seconds"

grep "$TEST_STR" ./${TEST_OUTPUT}

if [ $? -eq 0 ]
then
    echo "PASSED: the test log input was found in the output file"
else
    echo "FAILED: test log contents not present in results"
    echo "Here is the original:"
    cat ./${TEST_FILE}
    echo "Here are the results from the curl:"
    cat ./${TEST_OUTPUT}
    exit 1
fi

#rm -rf ${TEST_DIR}

#echo "Here is a comparison between the original log file with results from the GET:"
#diff ./${TEST_FILE} ./${TEST_OUTPUT} > diff.out
#cat diff.out
#echo ""

