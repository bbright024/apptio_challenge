#!/bin/bash

# Run this script to prepare the host OS & execute the logserver

# Enable/disable chrooting
NO_CHROOT=1

NEW_GROUP="apptiologserver"
SCP_DIR="./scpdir"
MAIN_APP_LOG_DIR="/home/bbright/mainapplog"
MAIN_APP_LOG_FILE="$MAIN_APP_LOG_DIR/mainapp.log"
LOG_USER="log-server-bbright"

# Changes group settings of the main app log and creates a new user
# that will run the logserver
groupadd $NEW_GROUP
mv $SCP_DIR/* $MAIN_APP_LOG_DIR
useradd $LOG_USER -g $NEW_GROUP
chown -R :$NEW_GROUP $MAIN_APP_LOG_DIR
chmod -R g+r $MAIN_APP_LOG_DIR

# Add cron rule to change group of files in the main app log directory,
#   in case new log files are created by the main app
crontab -l > tempcron
echo "00 00 * * * chown -R :$NEW_GROUP $MAIN_APP_LOG_DIR && chmod -R g+r $MAIN_APP_LOG_DIR " >> tempcron
crontab tempcron
rm tempcron


if [ $NO_CHROOT -eq 1 ]
then
    echo "chrooting not enabled.  exiting"
    exit 0
fi

###############################
##### UNDER CONSTRUCTION ######

# Now create a chroot jail in order to limit what an attacker can
# do after gaining shell access.  
cd  $MAIN_APP_LOG_DIR
mkdir bin dev etc etc/pam.d home lib lib/security var var/log usr usr/bin

# Copy needed libraries with ldd, awks, and find
# TODO: use find ... exec cp instead of xargs
ldd ./logserver | awk '{print $1}' | xargs -I{} cp "/lib/{}" /home/bbright/mainapplog/lib/

# Finally, run the chroot
chroot . /logserver ./conf.json
