#!/bin/bash

# This script tests the logserver's basic operability:
#   a) compile and run the logserver
#   b) use curl to read the test mainapp.log file
#   c) grep the results file for the test log entry


LOGSERVER_DIR="/home/bbright/go/src/apptio/logserver"
HOME_DIR="/home/bbright/go/src/apptio"

TEST_CONF="test_conf.json"
TEST_DIR="$HOME_DIR/test/"
TEST_TEMP="$TEST_DIR/test_temp"
TEST_FILE="mainapp.log"
TEST_STR="4/15/2018, This is a test string that should be downloaded "
TEST_SEARCH="This is a test string that should be downloaded"
TEST_OUTPUT="results.txt"
TIMEOUT=5


# First run the unit tests
cd $TEST_DIR
./small_test.sh
if [ ! $? -eq 0 ]
then
    echo "Error in small tests"
    exit 1
fi

# Now build the binary
cd $LOGSERVER_DIR
go build logserver.go

if [ ! -x ./logserver ]
then 
    echo "Error: logserver not found or not an executable file."
    exit
fi

cd $HOME_DIR

# Create temp directory
if [ ! -d ${TEST_TEMP} ]
then
    mkdir ${TEST_TEMP}
fi

# Copy needed files for logserver - conf.json, mainapp.log
cd $TEST_DIR
echo $TEST_STR > ./$TEST_FILE
echo "" >> ./$TEST_FILE
mv ./$TEST_FILE $TEST_TEMP/
mv $LOGSERVER_DIR/logserver $TEST_TEMP/
cp ./${TEST_CONF} ${TEST_TEMP}/
touch $TEST_TEMP/logserver.log

# Run the logserver
cd $TEST_TEMP
echo "Starting logserver with conf file:"
jq . ./$TEST_CONF
./logserver ${TEST_CONF} &
logserver_pid=$!

echo "Attempting to read mainapp.log from logserver:"

curl --max-time ${TIMEOUT} --silent --output ${TEST_OUTPUT} localhost:8888/read
(( $? == 28 )) && echo "Error: Fetch timed out after ${TIMEOUT} seconds"

# Kill the logserver, we're done with it
kill $logserver_pid
wait $logserver_pid

# Search the results for the message we saved in the mainapp.log entry
grep "$TEST_SEARCH" ./${TEST_OUTPUT}

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

# Clean up 
cd $TEST_DIR
rm -rf ${TEST_TEMP}

