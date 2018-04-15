#!/bin/bash

# Instead of having a long chain of
#[command] && \
#   [command] && \
# I decided to split the script in two - this script for the
# staging computer, and a second one to be run on the main app
# host machine

NEW_GROUP="apptiologserver"
ADMIN="bbright"
CONF_FILE="./host_machines_info.json"
SRC_CODE_DIR="/home/bbright/go/src/apptio/logserver"
SCP_DIR="./scpdir"
MAIN_APP_LOG_DIR="~/mainapplog/"
MAIN_APP_LOG_FILE="mainapp.log"
MAIN_LOG_FULL=$MAIN_APP_LOG_DIR$MAIN_APP_LOG_FILE
LOG_USER="logsserveruser"
# check if the conf file exists
if [ ! -f $CONF_FILE ]
then
    echo "no $CONF_FILE file found"
    exit 
fi

#MACHINES=$(jq .machines $CONF_FILE)
M_IP=$(jq -r .machines[0].ip $CONF_FILE )
M_ARCH=$(jq -r .machines[0].arch $CONF_FILE )
M_OS=$(jq -r .machines[0].os $CONF_FILE )

cd logserver
#GOARCH=$M_ARCH GOOS=$M_OS  go build logserver.go
go build logserver.go
cd ..

if [ -d $SCP_DIR ]
then
    rm -r $SCP_DIR
fi
   
mkdir $SCP_DIR
mv logserver/logserver $SCP_DIR
touch $SCP_DIR/logserver.log
cp ./install_on_host.sh $SCP_DIR

if [ "$M_IP" = "localhost" ]
then
    echo "this was a test"
    exit
fi



echo "copying files to $M_IP"
scp -r $SCP_DIR $ADMIN@$M_IP:~/
ssh -t $ADMIN@$M_IP "sudo ./scpdir/install_on_host.sh"

    


