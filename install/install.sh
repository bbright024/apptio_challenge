#!/bin/bash

# Basic script to install the logserver on a remote machine.
#    - Parses a json file with info for destination  machine
#    - compiles binary specific to host OS/architecture
#    - copies the binary and other installation files to dest.
#    - runs the configure.sh script on the destination to finish install


HOME_DIR="/home/bbright/go/src/apptio"
LOGSERVER_DIR="/home/bbright/go/src/apptio/logserver"
INSTALL_DIR="$HOME_DIR/install"

# Archive folder for storing hashes of binaries
ARCH_DIR="$HOME_DIR/archive"
# Staging folder to be copied to host
SCP_DIR="$HOME_DIR/install/scpdir"

# Name of user with root privileges on destination 
ADMIN="bbright"

# File containing destination host info
CONF_FILE="$INSTALL_DIR/host_info.json"

# File for saving sha512 hash of compiled binary
SHA_FILE="ls_sha_512.info"

TODAY=`date '+%Y_%m_%d_%H_%M_%S`
SHA_FILE_ARCHIVE="ls_sha_512_$TODAY"


# check if the file with host info exists
if [ ! -f $CONF_FILE ]
then
    echo "no $CONF_FILE file found"
    exit 
fi

# Parses host info file
# TODO:
#   write a while loop for every target - currently only installs
#   on first machine in array
# Also TODO: write a faster parsing mechanism than executing jq for every variable
M_IP=$(jq -r .machines[0].ip $CONF_FILE )
M_ARCH=$(jq -r .machines[0].arch $CONF_FILE )
M_OS=$(jq -r .machines[0].os $CONF_FILE )

# Create archive folder
if [ ! -d $ARCH_DIR ]
then
    mkdir $ARCH_DIR
fi

# Cleans up prior install files
if [ -d $SCP_DIR ]
then
    rm -rf $SCP_DIR
fi
mkdir $SCP_DIR

# Compile binary and store it's hash
cd $LOGSERVER_DIR
GOARCH=$M_ARCH GOOS=$M_OS  go build logserver.go
sha512sum logserver > ./${SHA_FILE}
cp -p ./${SHA_FILE} ./${SHA_FILE_ARCHIVE}
mv ./${SHA_FILE} ${SCP_DIR}/
mv ./${SHA_FILE_ARCHIVE} ${ARCH_DIR}/
cd $HOME_DIR


mv $LOGSERVER_DIR/logserver $SCP_DIR
touch $SCP_DIR/logserver.log
cp $HOME_DIR/configure.sh $SCP_DIR

# Cancels installation if dest host is localhost
if [ "$M_IP" = "localhost" ]
then
    echo "this was a test"
    exit
fi

# Copies files and runs configure script
echo "copying files to $M_IP"
scp -r $SCP_DIR $ADMIN@$M_IP:~/
ssh -t $ADMIN@$M_IP "sudo ./scpdir/configure.sh"





