#!/bin/bash


# This script must be run on the main app host machine

NEW_GROUP="apptiologserver"
SCP_DIR="./scpdir"
MAIN_APP_LOG_DIR="/home/bbright/mainapplog/"
MAIN_APP_LOG_FILE="mainapp.log"
MAIN_LOG_FULL=$MAIN_APP_LOG_DIR$MAIN_APP_LOG_FILE
LOG_USER="logsserveruser"

groupadd $NEW_GROUP
mv $SCP_DIR/* $MAIN_APP_LOG_DIR
useradd $LOG_USER -g $NEW_GROUP
chown -R :$NEW_GROUP $MAIN_APP_LOG_DIR
chmod -R g+r $MAIN_APP_LOG_DIR

# Now create a chroot environment
##### note:
##### heavily dependent on the main app host's environment
cd  $MAIN_APP_LOG_DIR
mkdir bin dev etc etc/pam.d home lib lib/security var var/log usr usr/bin
# Copy only the libraries needed to run chroot 


# Finally, run the chroot
chroot . /logserver
